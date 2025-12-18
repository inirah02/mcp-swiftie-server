// package mcpswiftieserver
package main

import (
	"encoding/json"
	"time"
)

// MCP Protocol Types
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

// Tool Types
type ToolInvocation struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

type ToolResult struct {
	Content interface{} `json:"content"`
	IsError bool        `json:"isError"`
}

// Taylor Swift Data Types
type Album struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	ReleaseYear int    `json:"release_year"`
	Era         string `json:"era"`
	Sales       int64  `json:"sales_millions"`
	Genre       string `json:"genre"`
}

type Song struct {
	ID         string `json:"id"`
	AlbumID    string `json:"album_id"`
	Title      string `json:"title"`
	Duration   int    `json:"duration_seconds"`
	Streams    int64  `json:"streams_millions"`
	ChartPeak  int    `json:"chart_peak"`
	GrammyNoms int    `json:"grammy_nominations"`
}

type Tour struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Year       int     `json:"year"`
	Shows      int     `json:"shows"`
	Attendance int64   `json:"attendance"`
	Revenue    float64 `json:"revenue_millions"`
}

type QueryResult struct {
	Columns   []string        `json:"columns"`
	Rows      [][]interface{} `json:"rows"`
	RowCount  int             `json:"row_count"`
	QueryTime time.Duration   `json:"query_time_ms"`
}
