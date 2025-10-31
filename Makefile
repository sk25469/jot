.PHONY: build install clean test

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

# Run with development mode
dev: build
	./jot

# Show version
version:
	@echo "jot v1.0.0-MVP"