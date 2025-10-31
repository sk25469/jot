.PHONY: build install clean test test-verbose test-coverage

# Build the binary
build:
	go build -o jot

# Install to system
install: build
	sudo mv jot /usr/local/bin/

# Install locally to ~/bin
install-user: build
	mkdir -p ~/bin
	mv jot ~/bin/
	@echo "Add ~/bin to your PATH if not already added"

# Clean build artifacts
clean:
	rm -f jot

# Run tests
test:
	go test ./...

# Run tests with verbose output
test-verbose:
	go test ./... -v

# Run tests with coverage report
test-coverage:
	go test ./... -cover

# Run tests with detailed coverage report
test-coverage-detailed:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run with development mode
dev: build
	./jot

# Show version
version:
	@echo "jot v1.0.0-MVP"

# Run all quality checks
check: test test-coverage
	go vet ./...
	go fmt ./...
	@echo "All checks passed!"