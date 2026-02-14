.PHONY: build run test vet lint clean

# Build the API server binary
build:
	go build -o bin/api ./cmd/api

# Run the API server
run:
	go run ./cmd/api

# Run all tests
test:
	go test -v -race -count=1 ./...

# Run go vet
vet:
	go vet ./...

# Clean build artifacts
clean:
	rm -rf bin/
