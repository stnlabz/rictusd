BINARY=bin/rictusd
CTL_BINARY=bin/rictusctl
SRC=./cmd/rictusd/main.go
CTL_SRC=./cmd/rictusctl/main.go

all: tidy build

tidy:
	@echo "🧹 Tidying modules..."
	@go mod tidy

build:
	@echo "🔨 Building RictusD Engine..."
	@mkdir -p bin
	@go mod tidy
	@go build -o $(BINARY) $(SRC)
	@echo "🔨 Building RictusD Controller..."
	@go build -o $(CTL_BINARY) $(CTL_SRC)

run:
	@./$(BINARY)

clean:
	@echo "🗑️ Purging binaries and caches..."
	@rm -rf bin/ data/threats.json
	@go clean -cache -modcache
	@echo "Purge complete."

help:
	@echo "RictusD Command Interface (Authorized):"
	@echo "  make build - Compile Engine and Controller"
	@echo "  make tidy  - Run go mod tidy"
	@echo "  make clean - Purge binaries and all Go caches"
	@echo "  make run   - Execute governed engine"
