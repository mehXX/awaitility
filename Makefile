lint:
	@echo "Running golangci-lint..."
	golangci-lint run

t:
	@echo "Running tests"
	go test ./... -race