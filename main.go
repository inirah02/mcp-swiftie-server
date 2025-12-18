package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	// Metrics
	queriesExecuted  atomic.Int64
	totalLatency     atomic.Int64
	activeGoroutines atomic.Int32
)

type Metrics struct {
	QueriesExecuted  int64   `json:"queries_executed"`
	AvgLatencyMS     float64 `json:"avg_latency_ms"`
	ActiveGoroutines int32   `json:"active_goroutines"`
	UptimeSeconds    int64   `json:"uptime_seconds"`
}

var startTime time.Time

func main() {
	startTime = time.Now()

	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
	log.Println("[INFO] 🎤 MCP Swiftie Server starting...")

	server := NewServer()

	// Register tools
	tools := server.ListTools()
	log.Printf("[INFO] Registered %d tools: %v", len(tools), getToolNames(tools))

	// HTTP handlers
	http.HandleFunc("/mcp", func(w http.ResponseWriter, r *http.Request) {
		handleMCPConnection(w, r, server)
	})

	http.HandleFunc("/metrics", handleMetrics)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}

	addr := fmt.Sprintf(":%s", port)
	log.Printf("[INFO] Server listening on %s", addr)
	log.Println("[INFO] Ready for connections ✨")

	// Graceful shutdown
	srv := &http.Server{Addr: addr}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[ERROR] Server failed: %v", err)
		}
	}()

	// Wait for interrupt
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("[INFO] Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("[ERROR] Server forced to shutdown: %v", err)
	}

	log.Println("[INFO] Server exited")
}

func handleMCPConnection(w http.ResponseWriter, r *http.Request, server *Server) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[ERROR] WebSocket upgrade failed: %v", err)
		return
	}
	defer conn.Close()

	log.Printf("[INFO] New MCP connection from %s", r.RemoteAddr)

	// Send server info
	serverInfo := MCPResponse{
		JSONRPC: "2.0",
		ID:      uuid.New().String(),
		Result: map[string]interface{}{
			"protocolVersion": "0.1.0",
			"serverInfo": map[string]string{
				"name":    "mcp-swiftie-server",
				"version": "1.0.0",
			},
			"capabilities": map[string]interface{}{
				"tools": map[string]interface{}{},
			},
		},
	}

	if err := conn.WriteJSON(serverInfo); err != nil {
		log.Printf("[ERROR] Failed to send server info: %v", err)
		return
	}

	// Handle requests
	for {
		var req MCPRequest
		if err := conn.ReadJSON(&req); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("[ERROR] WebSocket error: %v", err)
			}
			break
		}

		go handleMCPRequest(conn, req, server)
	}

	log.Printf("[INFO] Connection closed from %s", r.RemoteAddr)
}

func handleMCPRequest(conn *websocket.Conn, req MCPRequest, server *Server) {
	activeGoroutines.Add(1)
	defer activeGoroutines.Add(-1)

	start := time.Now()

	var response MCPResponse
	response.JSONRPC = "2.0"
	response.ID = req.ID

	switch req.Method {
	case "tools/list":
		response.Result = map[string]interface{}{
			"tools": server.ListTools(),
		}

	case "tools/call":
		var invocation ToolInvocation
		if err := json.Unmarshal(req.Params, &invocation); err != nil {
			response.Error = &MCPError{Code: -32600, Message: "Invalid params"}
			break
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		result := server.ExecuteTool(ctx, invocation)

		if result.IsError {
			response.Error = &MCPError{Code: -32000, Message: result.Content.(string)}
		} else {
			response.Result = result.Content
		}

		// Update metrics
		queriesExecuted.Add(1)
		totalLatency.Add(time.Since(start).Milliseconds())

	default:
		response.Error = &MCPError{Code: -32601, Message: "Method not found"}
	}

	if err := conn.WriteJSON(response); err != nil {
		log.Printf("[ERROR] Failed to send response: %v", err)
	}
}

func handleMetrics(w http.ResponseWriter, r *http.Request) {
	queries := queriesExecuted.Load()
	latency := totalLatency.Load()

	avgLatency := float64(0)
	if queries > 0 {
		avgLatency = float64(latency) / float64(queries)
	}

	metrics := Metrics{
		QueriesExecuted:  queries,
		AvgLatencyMS:     avgLatency,
		ActiveGoroutines: activeGoroutines.Load(),
		UptimeSeconds:    int64(time.Since(startTime).Seconds()),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}

func getToolNames(tools []map[string]interface{}) []string {
	names := make([]string, len(tools))
	for i, tool := range tools {
		names[i] = tool["name"].(string)
	}
	return names
}
