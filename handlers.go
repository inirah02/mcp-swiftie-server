package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

type Server struct {
	presto *PrestoClient
}

func NewServer() *Server {
	return &Server{
		presto: NewPrestoClient(),
	}
}

// ListTools returns available MCP tools
func (s *Server) ListTools() []map[string]interface{} {
	return []map[string]interface{}{
		{
			"name":        "list_tables",
			"description": "List all available tables in the Taylor Swift database",
			"inputSchema": map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
		},
		{
			"name":        "query_albums",
			"description": "Query Taylor Swift albums with optional filters",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"era": map[string]string{
						"type":        "string",
						"description": "Filter by era (e.g., 'Pop', 'Country', 'Indie Folk')",
					},
				},
			},
		},
		{
			"name":        "query_songs",
			"description": "Query Taylor Swift songs with streaming and chart data",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"album_id": map[string]string{
						"type":        "string",
						"description": "Filter by album ID",
					},
					"min_streams": map[string]interface{}{
						"type":        "number",
						"description": "Minimum streams in millions",
					},
				},
			},
		},
		{
			"name":        "analyze_tours",
			"description": "Analyze Taylor Swift tour data including revenue and attendance",
			"inputSchema": map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
		},
		{
			"name":        "streaming_query",
			"description": "Execute a large query with streaming results (for demo)",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"table": map[string]string{
						"type":        "string",
						"description": "Table to query (albums, songs, tours)",
					},
				},
				"required": []string{"table"},
			},
		},
	}
}

// ExecuteTool handles tool invocation
func (s *Server) ExecuteTool(ctx context.Context, invocation ToolInvocation) ToolResult {
	log.Printf("[INFO] Tool invocation: %s", invocation.Name)
	log.Printf("[DEBUG] Arguments: %v", invocation.Arguments)

	switch invocation.Name {
	case "list_tables":
		return s.handleListTables(ctx)
	case "query_albums":
		return s.handleQueryAlbums(ctx, invocation.Arguments)
	case "query_songs":
		return s.handleQuerySongs(ctx, invocation.Arguments)
	case "analyze_tours":
		return s.handleAnalyzeTours(ctx)
	case "streaming_query":
		return s.handleStreamingQuery(ctx, invocation.Arguments)
	default:
		return ToolResult{
			Content: fmt.Sprintf("Unknown tool: %s", invocation.Name),
			IsError: true,
		}
	}
}

func (s *Server) handleListTables(ctx context.Context) ToolResult {
	start := time.Now()

	result, err := s.presto.Query(ctx, "SHOW TABLES")
	if err != nil {
		return ToolResult{Content: err.Error(), IsError: true}
	}

	log.Printf("[INFO] Tool completed in %v", time.Since(start))
	return ToolResult{Content: result, IsError: false}
}

func (s *Server) handleQueryAlbums(ctx context.Context, args map[string]interface{}) ToolResult {
	start := time.Now()

	sql := "SELECT * FROM albums"
	result, err := s.presto.Query(ctx, sql)
	if err != nil {
		return ToolResult{Content: err.Error(), IsError: true}
	}

	log.Printf("[INFO] Returned %d rows in %v", result.RowCount, time.Since(start))
	return ToolResult{Content: result, IsError: false}
}

func (s *Server) handleQuerySongs(ctx context.Context, args map[string]interface{}) ToolResult {
	start := time.Now()

	sql := "SELECT * FROM songs"
	result, err := s.presto.Query(ctx, sql)
	if err != nil {
		return ToolResult{Content: err.Error(), IsError: true}
	}

	log.Printf("[INFO] Returned %d rows in %v", result.RowCount, time.Since(start))
	return ToolResult{Content: result, IsError: false}
}

func (s *Server) handleAnalyzeTours(ctx context.Context) ToolResult {
	start := time.Now()

	sql := "SELECT * FROM tours"
	result, err := s.presto.Query(ctx, sql)
	if err != nil {
		return ToolResult{Content: err.Error(), IsError: true}
	}

	log.Printf("[INFO] Returned %d rows in %v", result.RowCount, time.Since(start))
	return ToolResult{Content: result, IsError: false}
}

func (s *Server) handleStreamingQuery(ctx context.Context, args map[string]interface{}) ToolResult {
	start := time.Now()
	table := args["table"].(string)

	sql := fmt.Sprintf("SELECT * FROM %s", table)

	// Use streaming with batches
	rowsChan, errChan := s.presto.StreamQuery(ctx, sql, 5)

	batchCount := 0
	totalRows := 0

	for {
		select {
		case batch, ok := <-rowsChan:
			if !ok {
				log.Printf("[INFO] Streaming completed: %d batches, %d rows in %v",
					batchCount, totalRows, time.Since(start))
				return ToolResult{
					Content: map[string]interface{}{
						"batches":    batchCount,
						"total_rows": totalRows,
						"query_time": time.Since(start).Milliseconds(),
					},
					IsError: false,
				}
			}
			batchCount++
			totalRows += len(batch)
			log.Printf("[DEBUG] Streaming batch %d (%d rows)", batchCount, len(batch))

		case err := <-errChan:
			if err != nil {
				return ToolResult{Content: err.Error(), IsError: true}
			}

		case <-ctx.Done():
			log.Printf("[WARN] Context cancelled: %v", ctx.Err())
			return ToolResult{Content: "Query cancelled", IsError: true}
		}
	}
}

// ExecuteToolsConcurrently demonstrates parallel tool execution
func (s *Server) ExecuteToolsConcurrently(ctx context.Context, tools []ToolInvocation) []ToolResult {
	results := make(chan ToolResult, len(tools))

	for _, tool := range tools {
		go func(t ToolInvocation) {
			// Add timeout per tool
			toolCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
			defer cancel()

			result := s.ExecuteTool(toolCtx, t)
			results <- result
		}(tool)
	}

	// Collect results
	var output []ToolResult
	for range tools {
		output = append(output, <-results)
	}

	return output
}
