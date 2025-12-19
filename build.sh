#!/bin/bash

set -e

echo "Building MCP Swiftie Server..."

# Clean
rm -f mcp-server mcp-client

# Build server
echo "Building server..."
go build -o mcp-server main.go handlers.go presto.go types.go

# Build client from examples
echo "Building client..."
go build -o mcp-client examples/simple_client.go

echo ""
echo "Build successful!"
ls -lh mcp-server mcp-client

echo ""
echo "To run:"
echo "   Terminal 1: ./mcp-server"
echo "   Terminal 2: ./mcp-client"
echo "   Terminal 3: go test -bench=. -benchmem"
