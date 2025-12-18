package main

import (
	"context"
	"fmt"
	"sync"
	"testing"
)

// Benchmark single query execution
func BenchmarkSingleQuery(b *testing.B) {
	server := NewServer()
	ctx := context.Background()

	invocation := ToolInvocation{
		Name:      "query_albums",
		Arguments: map[string]interface{}{},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		server.ExecuteTool(ctx, invocation)
	}
}

// Benchmark concurrent queries
func BenchmarkConcurrentQueries(b *testing.B) {
	server := NewServer()

	concurrencyLevels := []int{10, 50, 100}

	for _, concurrency := range concurrencyLevels {
		b.Run(fmt.Sprintf("Concurrency-%d", concurrency), func(b *testing.B) {
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				var wg sync.WaitGroup
				ctx := context.Background()

				for j := 0; j < concurrency; j++ {
					wg.Add(1)
					go func() {
						defer wg.Done()

						invocation := ToolInvocation{
							Name:      "query_songs",
							Arguments: map[string]interface{}{},
						}

						server.ExecuteTool(ctx, invocation)
					}()
				}

				wg.Wait()
			}
		})
	}
}

// Benchmark streaming queries
func BenchmarkStreamingQuery(b *testing.B) {
	server := NewServer()
	ctx := context.Background()

	invocation := ToolInvocation{
		Name: "streaming_query",
		Arguments: map[string]interface{}{
			"table": "songs",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		server.ExecuteTool(ctx, invocation)
	}
}

// Benchmark tool listing
func BenchmarkListTools(b *testing.B) {
	server := NewServer()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		server.ListTools()
	}
}
