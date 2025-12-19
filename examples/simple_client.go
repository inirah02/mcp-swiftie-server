package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type MCPRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      string          `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type MCPResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      string      `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *MCPError   `json:"error,omitempty"`
}

type MCPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ToolInvocation struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

func main() {
	// Connect to server
	url := "ws://localhost:9000/mcp"
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	log.Println("Connected to MCP Swiftie Server")

	// Read server info
	var serverInfo MCPResponse
	if err := conn.ReadJSON(&serverInfo); err != nil {
		log.Fatalf("Failed to read server info: %v", err)
	}
	log.Printf("Server info: %+v\n", serverInfo.Result)

	// List tools
	log.Println("\n Listing available tools...")
	listToolsReq := MCPRequest{
		JSONRPC: "2.0",
		ID:      uuid.New().String(),
		Method:  "tools/list",
	}

	if err := conn.WriteJSON(listToolsReq); err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}

	var toolsResp MCPResponse
	if err := conn.ReadJSON(&toolsResp); err != nil {
		log.Fatalf("Failed to read response: %v", err)
	}

	result := toolsResp.Result.(map[string]interface{})
	tools := result["tools"].([]interface{})
	log.Printf("Available tools: %d\n", len(tools))
	for _, t := range tools {
		tool := t.(map[string]interface{})
		log.Printf("  - %s: %s", tool["name"], tool["description"])
	}

	// Query albums
	log.Println("\n Querying Taylor Swift albums...")
	queryReq := MCPRequest{
		JSONRPC: "2.0",
		ID:      uuid.New().String(),
		Method:  "tools/call",
	}

	invocation := ToolInvocation{
		Name:      "query_albums",
		Arguments: map[string]interface{}{},
	}

	params, _ := json.Marshal(invocation)
	queryReq.Params = params

	start := time.Now()
	if err := conn.WriteJSON(queryReq); err != nil {
		log.Fatalf("Failed to send query: %v", err)
	}

	var queryResp MCPResponse
	if err := conn.ReadJSON(&queryResp); err != nil {
		log.Fatalf("Failed to read response: %v", err)
	}
	duration := time.Since(start)

	if queryResp.Error != nil {
		log.Fatalf("Query failed: %v", queryResp.Error)
	}

	queryResult := queryResp.Result.(map[string]interface{})
	log.Printf("Query completed in %v", duration)
	log.Printf("Result: %+v\n", queryResult)

	// Query songs with streaming
	log.Println("\n🎵 Testing streaming query...")
	streamReq := MCPRequest{
		JSONRPC: "2.0",
		ID:      uuid.New().String(),
		Method:  "tools/call",
	}

	streamInvocation := ToolInvocation{
		Name: "streaming_query",
		Arguments: map[string]interface{}{
			"table": "songs",
		},
	}

	streamParams, _ := json.Marshal(streamInvocation)
	streamReq.Params = streamParams

	start = time.Now()
	if err := conn.WriteJSON(streamReq); err != nil {
		log.Fatalf("Failed to send streaming query: %v", err)
	}

	var streamResp MCPResponse
	if err := conn.ReadJSON(&streamResp); err != nil {
		log.Fatalf("Failed to read response: %v", err)
	}
	duration = time.Since(start)

	log.Printf("Streaming query completed in %v", duration)
	log.Printf("Result: %+v\n", streamResp.Result)

	// Analyze tours
	log.Println("\n🎤 Analyzing tour data...")
	tourReq := MCPRequest{
		JSONRPC: "2.0",
		ID:      uuid.New().String(),
		Method:  "tools/call",
	}

	tourInvocation := ToolInvocation{
		Name:      "analyze_tours",
		Arguments: map[string]interface{}{},
	}

	tourParams, _ := json.Marshal(tourInvocation)
	tourReq.Params = tourParams

	start = time.Now()
	if err := conn.WriteJSON(tourReq); err != nil {
		log.Fatalf("Failed to send tour query: %v", err)
	}

	var tourResp MCPResponse
	if err := conn.ReadJSON(&tourResp); err != nil {
		log.Fatalf("Failed to read response: %v", err)
	}
	duration = time.Since(start)

	log.Printf("Tour analysis completed in %v", duration)

	tourResult := tourResp.Result.(map[string]interface{})
	rows := tourResult["rows"].([]interface{})

	log.Println("\nTour Revenue Summary:")
	for _, row := range rows {
		r := row.([]interface{})
		tourName := r[1].(string)
		revenue := r[5].(float64)
		log.Printf("  %s: $%.1fM", tourName, revenue)
	}

	log.Println("\n✨ Demo completed successfully!")
}
