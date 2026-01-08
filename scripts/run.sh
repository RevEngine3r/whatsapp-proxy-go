#!/bin/bash
# WhatsApp Proxy Go - Unix Runner Script
# Usage: ./run.sh [config_file]

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
CONFIG_FILE="${1:-${CONFIG_FILE:-config.yaml}}"
LOG_FILE="${LOG_FILE:-whatsapp-proxy.log}"
BINARY="./whatsapp-proxy"

# Check if binary exists
if [ ! -f "$BINARY" ]; then
    # Try dist directory
    if [ -f "./dist/whatsapp-proxy" ]; then
        BINARY="./dist/whatsapp-proxy"
    else
        echo -e "${RED}Error: Binary not found${NC}"
        echo "Please build the project first: make build"
        exit 1
    fi
fi

# Check if config file exists
if [ ! -f "$CONFIG_FILE" ]; then
    echo -e "${RED}Error: Config file not found: $CONFIG_FILE${NC}"
    echo "Usage: $0 [config_file]"
    echo ""
    echo "You can also set CONFIG_FILE environment variable:"
    echo "  export CONFIG_FILE=config.yaml"
    echo ""
    echo "Copy the example config to get started:"
    echo "  cp configs/config.example.yaml config.yaml"
    exit 1
fi

# Make binary executable
chmod +x "$BINARY"

echo -e "${GREEN}Starting WhatsApp Proxy Go${NC}"
echo -e "${YELLOW}Config:${NC} $CONFIG_FILE"
echo -e "${YELLOW}Logs:${NC} $LOG_FILE"
echo ""

# Run the proxy with logging
"$BINARY" --config "$CONFIG_FILE" 2>&1 | tee -a "$LOG_FILE"
