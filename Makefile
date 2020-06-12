.PHONY: fmt
fmt:
	@echo "==> Running gofmt..."
	gofmt -s -w .

.PHONY: build
build: fmt test
	@echo "==> Building library..."
	go build -ldflags="-s -w" ./...
	@echo "==> Building the CLI..."
	go build -ldflags="-s -w" ./cmd/whathappens

.PHONY: test
test:
	@echo "==> Running tests..."
	@go test -cover .
