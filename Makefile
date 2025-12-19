.PHONY: all build test bench clean

all: build

build:
	@./build.sh

test:
	@echo "Running tests..."
	@go test -v -race

bench:
	@echo "📊 Running benchmarks..."
	@go test -bench=. -benchmem | tee benchmark_results.txt
	@echo ""
	@echo "Benchmark results saved!"

clean:
	@rm -f mcp-server mcp-client benchmark_results.txt
