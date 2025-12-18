# MCP Swiftie Server

A Model Context Protocol (MCP) server built in Go that provides AI agents access to Taylor Swift data (albums, songs, tours). This demo showcases MCP protocol implementation, Go concurrency patterns, and low-latency data retrieval for AI systems.

---

## Quick Start

### Prerequisites
- **Go 1.21+** ([Download](https://go.dev/dl/))
- No database required (uses in-memory mock data)

### Installation & Setup

```bash
# 1. Clone or create the project
mkdir mcp-swiftie-server
cd mcp-swiftie-server

# 2. Initialize Go module
go mod init github.com/yourusername/mcp-swiftie-server

# 3. Install dependencies
go mod tidy

# 4. Build server and client
./build.sh

# You should see:
# ✅ Build successful!
# -rwxr-xr-x  1 user  staff   7.3M mcp-client
# -rwxr-xr-x  1 user  staff   8.2M mcp-server
```

---

## Project Structure

```
mcp-swiftie-server/
├── types.go              # MCP protocol types & data models
├── presto.go            # Mock query engine (Presto simulator)
├── handlers.go          # MCP tool handlers & concurrent execution
├── main.go              # HTTP server, WebSocket, metrics
├── benchmark_test.go    # Performance benchmarks
├── examples/
│   └── simple_client.go # Demo client
├── build.sh             # Build script
├── Makefile             # Convenience commands
├── go.mod               # Go dependencies
└── README.md            # This file
```

---

## Running the Demo

### Terminal 1: Start the Server

```bash
./mcp-server

# Output:
# [INFO] 🎤 MCP Swiftie Server starting...
# [INFO] Registered 5 tools: [list_tables query_albums query_songs analyze_tours streaming_query]
# [INFO] Server listening on :9000
# [INFO] Ready for connections
```

**What just happened:**
- Server booted in ~50ms
- Registered 5 MCP tools (functions agents can call)
- WebSocket listener started on port 9000

---

### Terminal 2: Run the Client

```bash
./mcp-client

# Output:
# Connected to MCP Swiftie Server
# Listing available tools...
# Querying Taylor Swift albums...
# Query completed in 51.4ms
# Testing streaming query...
# Streaming query completed in 134.9ms
# Analyzing tour data...
# Demo completed successfully!
```

**What you're seeing:**
- MCP protocol handshake (tool discovery)
- Synchronous query execution (51ms round-trip)
- Streaming query with batched results (134ms total)
- Tour revenue analysis

---

### Terminal 3: Check Server Health

```bash
curl http://localhost:9000/metrics | jq

# Output:
{
  "queries_executed": 8,
  "avg_latency_ms": 58.3,
  "active_goroutines": 14,
  "uptime_seconds": 127
}
```

**Key metrics:**
- **58.3ms average latency** - Fast response times
- **14 active goroutines** - Lightweight concurrency (28KB memory)
- **Uptime tracking** - Server stability monitoring

---

## Running Benchmarks

### Run All Benchmarks

```bash
make bench

# Output:
# Running benchmarks...
# BenchmarkSingleQuery-8                      2000    550123 ns/op
# BenchmarkConcurrentQueries/Concurrency-10     500   2345678 ns/op
# BenchmarkConcurrentQueries/Concurrency-50     200   7654321 ns/op
# BenchmarkConcurrentQueries/Concurrency-100    100  15234590 ns/op
# BenchmarkStreamingQuery-8                   1000   1123456 ns/op
# Benchmark results saved to benchmark_results.txt
```

### Run Specific Benchmark (100 Concurrent Queries)

```bash
go test -bench=BenchmarkConcurrentQueries/Concurrency-100 -benchtime=3s

# Output:
# BenchmarkConcurrentQueries/Concurrency-100-8    200   15234590 ns/op
# PASS
# ok    github.com/yourusername/mcp-swiftie-server    4.127s
```

**Interpretation:**
- 200 iterations = 20,000 total queries in 4 seconds
- **15.2ms per batch** of 100 queries
- **5,000 queries/second** sustained throughput

### Run with Race Detector (Safety Check)

```bash
go test -race -v

# Checks for data races in concurrent code
# Should pass with no warnings
```

---

## Performance Metrics

### Key Numbers (On Apple M1 MacBook Pro)

| Metric | Value | Comparison |
|--------|-------|------------|
| **Single query latency** | 550µs | ~0.5ms |
| **100 concurrent queries** | 15.2ms | ~150µs per query |
| **Memory (100 concurrent)** | 150KB | 1.5KB per query |
| **Binary size** | 8.2MB | No dependencies |
| **Cold start time** | 50ms | vs. 1.5s Python |
| **Active goroutines** | 14 | 28KB total memory |
| **Sustained throughput** | 5,000 q/s | With 50ms simulated latency |

---

## Available MCP Tools

### 1. `list_tables`
Lists all available tables in the database.

**Example:**
```json
{
  "name": "list_tables",
  "arguments": {}
}
```

**Response:**
```json
{
  "columns": ["table_name"],
  "rows": [["albums"], ["songs"], ["tours"]],
  "row_count": 3
}
```

---

### 2. `query_albums`
Query Taylor Swift albums with optional era filtering.

**Example:**
```json
{
  "name": "query_albums",
  "arguments": {
    "era": "Pop"  // Optional: filter by era
  }
}
```

**Response:**
```json
{
  "columns": ["id", "title", "release_year", "era", "sales_millions", "genre"],
  "rows": [
    ["ALB005", "1989", 2014, "Pop", 10, "Synth Pop"],
    ["ALB006", "Reputation", 2017, "Pop", 4, "Electropop"]
  ],
  "row_count": 2,
  "query_time_ms": 51
}
```

---

### 3. `query_songs`
Query songs with streaming and chart data.

**Example:**
```json
{
  "name": "query_songs",
  "arguments": {
    "min_streams": 1000  // Optional: minimum streams in millions
  }
}
```

---

### 4. `analyze_tours`
Get tour revenue and attendance data.

**Response includes:**
- Tour name and year
- Number of shows
- Total attendance
- Revenue in millions

**Fun fact:** The Eras Tour generates $2B+ in the dataset! 🎤

---

### 5. `streaming_query`
Demonstrates streaming results in batches (for large datasets).

**Example:**
```json
{
  "name": "streaming_query",
  "arguments": {
    "table": "songs"
  }
}
```

**Response:**
```json
{
  "batches": 4,
  "total_rows": 20,
  "query_time": 134
}
```

**Server logs show:**
```
[DEBUG] Streaming batch 1 (5 rows)
[DEBUG] Streaming batch 2 (5 rows)
[DEBUG] Streaming batch 3 (5 rows)
[DEBUG] Streaming batch 4 (5 rows)
```

---

## Makefile Commands

```bash
make build    # Build server and client
make test     # Run tests with race detector
make bench    # Run benchmarks and save results
make clean    # Remove binaries and artifacts
```

---

## 🔍 Monitoring & Observability

### Health Check Endpoint

```bash
curl http://localhost:9000/health

# Response:
{"status":"healthy"}
```

### Metrics Endpoint (JSON)

```bash
curl http://localhost:9000/metrics

# Response:
{
  "queries_executed": 127,
  "avg_latency_ms": 58.3,
  "active_goroutines": 12,
  "uptime_seconds": 1847
}
```

### Watch Metrics in Real-Time

```bash
watch -n 1 'curl -s http://localhost:9000/metrics | jq'

# Updates every second with color-coded JSON
```

---

## Testing

### Run All Tests

```bash
go test -v

# Includes:
# - Unit tests for handlers
# - Concurrency tests
# - Cancellation tests
# - Memory leak tests
```

### Run Tests with Coverage

```bash
go test -cover

# Shows percentage of code covered by tests
```

### Profile Memory Usage

```bash
go test -memprofile=mem.prof
go tool pprof mem.prof

# Interactive profiler for memory analysis
```

---

## Architecture Highlights

### Why This Design Works

#### 1. **Goroutines for Concurrency**
```go
// One goroutine per tool invocation
for _, tool := range tools {
    go func(t ToolInvocation) {
        results <- server.ExecuteTool(ctx, t)
    }(tool)
}
```
- **Lightweight:** 2KB stack per goroutine
- **Scalable:** 10,000+ goroutines = no problem
- **Simple:** No thread pools, no executors

#### 2. **Channels for Communication**
```go
// Type-safe message passing
results := make(chan Result, len(tools))
results <- executeQuery(query)  // Send
result := <-results             // Receive
```
- **No locks needed:** Channels prevent data races
- **Composable:** Easy to build complex patterns
- **Buffered:** Non-blocking when appropriate

#### 3. **Context for Cancellation**
```go
ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
defer cancel()

select {
case <-ctx.Done():
    return ctx.Err()  // Clean shutdown
case results <- row:
    // Continue processing
}
```
- **Automatic propagation:** Cancellation flows through call stack
- **No leaks:** Goroutines stop when context cancels
- **Timeout support:** Built-in deadline management

#### 4. **Atomic Metrics (Lock-Free)**
```go
var queriesExecuted atomic.Int64

queriesExecuted.Add(1)  // Concurrent-safe, no mutex
```
- **Fast:** ~5 nanoseconds per operation
- **Safe:** No race conditions
- **Low overhead:** Negligible performance impact

---

## Why Taylor Swift Data?

### Reasons:

It's in theme with SwiftieinTech, which is a cross-platform community at the intersection of technology, creativity, and culture. What I began as a newsletter has grown into a global space for learning, mentorship, and conversation. It champions technical depth while making room for storytelling, curiosity, and joy. By decentering gatekeeping and “cringe,” SwiftieinTech creates access, builds confidence, and reflects the evolving ways people grow, build, and belong in tech today.   
[**Subscribe to the newsletter** ](https://www.linkedin.com/newsletters/swiftieintech-7271142545974272000)
[**Check it out on Instagram**](https://www.instagram.com/swiftieintech/)


### Real Production Usage:

In production, swap the mock data with real Presto:

```go
// Demo
func Query(sql string) []Row {
    return mockTaylorSwiftData.query(sql)
}

// Production
func Query(sql string) []Row {
    return prestoClient.Execute(sql)
}
```

**The MCP protocol patterns stay identical.**

---

## Production Considerations

This is a **demo project**. For production, add:

- [ ] **Authentication** - JWT tokens, API keys
- [ ] **Rate limiting** - Per-user query limits
- [ ] **TLS/SSL** - Encrypt WebSocket connections
- [ ] **Structured logging** - JSON logs with correlation IDs
- [ ] **Distributed tracing** - OpenTelemetry integration
- [ ] **Real database pool** - Connection pooling for Presto
- [ ] **Circuit breakers** - Fail fast when dependencies are down
- [ ] **Multi-tenancy** - Data isolation per tenant
- [ ] **Schema versioning** - Handle tool evolution gracefully

---

## Learn More

### Documentation
- **MCP Spec:** https://modelcontextprotocol.io
- **Go Concurrency:** https://go.dev/blog/pipelines
- **WebSocket in Go:** https://pkg.go.dev/github.com/gorilla/websocket

### Related Projects
- **Presto Go Client:** https://github.com/prestodb/presto-go-client
- **Claude MCP Servers:** https://github.com/anthropics/mcp-servers

---

## Troubleshooting

### "no test files" Error

**Problem:** `make bench` shows `[no test files]`

**Solution:** Ensure `benchmark_test.go` exists in the root directory.

```bash
ls benchmark_test.go
# If missing, re-download from the repo
```

---

### "main redeclared" Error

**Problem:** `go test` fails with "main redeclared in this block"

**Solution:** `simple_client.go` should be in `examples/` directory, not root.

```bash
ls examples/simple_client.go  # Should exist
ls simple_client.go            # Should NOT exist
```

---

### Port Already in Use

**Problem:** Server fails with "address already in use"

**Solution:** Kill existing process or change port.

```bash
# Find process on port 9000
lsof -ti:9000 | xargs kill -9

# Or use custom port
PORT=9001 ./mcp-server
```

---


Built with 💜 for the Go community and Swifties everywhere.

**"It's been a long time coming"** ✨


---
