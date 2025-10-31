#!/bin/bash

# jot Quick Install Script
# Usage: curl -sSL https://raw.githubusercontent.com/sahilsarwar/jot/main/install.sh | bash

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $ARCH in
    x86_64) ARCH="amd64" ;;
    aarch64) ARCH="arm64" ;;
    armv7l) ARCH="arm" ;;
    *) echo -e "${RED}Unsupported architecture: $ARCH${NC}"; exit 1 ;;
esac

echo -e "${GREEN}Installing jot for $OS/$ARCH...${NC}"

# Check if Go is installed for building from source
if command -v go &> /dev/null; then
    echo -e "${YELLOW}Go detected. Building from source...${NC}"
    
    # Create temporary directory
    TMP_DIR=$(mktemp -d)
    cd $TMP_DIR
    
    # Clone and build
    git clone https://github.com/sk25469/jot.git
    cd jot
    go build -o jot
    
    # Install
    if [[ $EUID -eq 0 ]]; then
        mv jot /usr/local/bin/
        echo -e "${GREEN}jot installed to /usr/local/bin/${NC}"
    else
        mkdir -p ~/bin
        mv jot ~/bin/
        echo -e "${GREEN}jot installed to ~/bin/${NC}"
        echo -e "${YELLOW}Make sure ~/bin is in your PATH${NC}"
    fi
    
    # Cleanup
    cd /
    rm -rf $TMP_DIR
    
else
    echo -e "${RED}Go not found. Please install Go or download a pre-built binary.${NC}"
    echo -e "${YELLOW}Visit: https://golang.org/dl/${NC}"
    exit 1
fi

echo -e "${GREEN}Installation complete!${NC}"
echo -e "Try: ${YELLOW}jot new \"My first note\"${NC}"