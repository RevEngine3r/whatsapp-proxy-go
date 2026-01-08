#!/bin/bash
# WhatsApp Proxy Go - Linux Service Installation Script
# Installs the proxy as a systemd service

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Check if running as root
if [ "$EUID" -ne 0 ]; then
    echo -e "${RED}Error: This script must be run as root${NC}"
    echo "Usage: sudo $0"
    exit 1
fi

# Configuration
SERVICE_NAME="whatsapp-proxy"
USER="whatsapp-proxy"
GROUP="whatsapp-proxy"
INSTALL_DIR="/opt/whatsapp-proxy"
CONFIG_DIR="/etc/whatsapp-proxy"
BINARY_NAME="whatsapp-proxy"

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  WhatsApp Proxy - Service Installer${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Check if binary exists
if [ ! -f "$BINARY_NAME" ] && [ ! -f "dist/$BINARY_NAME" ]; then
    echo -e "${RED}Error: Binary not found${NC}"
    echo "Please build the project first: make build"
    exit 1
fi

# Use dist binary if main not found
if [ ! -f "$BINARY_NAME" ] && [ -f "dist/$BINARY_NAME" ]; then
    BINARY_NAME="dist/$BINARY_NAME"
fi

# Check if systemd service file exists
if [ ! -f "systemd/whatsapp-proxy.service" ]; then
    echo -e "${RED}Error: systemd service file not found${NC}"
    exit 1
fi

# Create service user
echo -e "${YELLOW}Creating service user...${NC}"
if ! id -u "$USER" >/dev/null 2>&1; then
    useradd -r -s /bin/false -d "$INSTALL_DIR" "$USER"
    echo -e "${GREEN}  ✓ User created: $USER${NC}"
else
    echo -e "${YELLOW}  → User already exists: $USER${NC}"
fi

# Create directories
echo -e "${YELLOW}Creating directories...${NC}"
mkdir -p "$INSTALL_DIR"
mkdir -p "$CONFIG_DIR"
echo -e "${GREEN}  ✓ Install directory: $INSTALL_DIR${NC}"
echo -e "${GREEN}  ✓ Config directory: $CONFIG_DIR${NC}"

# Copy binary
echo -e "${YELLOW}Installing binary...${NC}"
cp "$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
chmod +x "$INSTALL_DIR/$BINARY_NAME"
echo -e "${GREEN}  ✓ Binary installed${NC}"

# Copy or create config
if [ -f "config.yaml" ]; then
    echo -e "${YELLOW}Installing existing config...${NC}"
    cp config.yaml "$CONFIG_DIR/config.yaml"
    echo -e "${GREEN}  ✓ Config installed from config.yaml${NC}"
elif [ -f "configs/config.example.yaml" ]; then
    echo -e "${YELLOW}Installing example config...${NC}"
    cp configs/config.example.yaml "$CONFIG_DIR/config.yaml"
    echo -e "${YELLOW}  → Please edit $CONFIG_DIR/config.yaml${NC}"
else
    echo -e "${RED}Warning: No config file found${NC}"
    echo -e "${YELLOW}  → Please create $CONFIG_DIR/config.yaml${NC}"
fi

# Set ownership
echo -e "${YELLOW}Setting permissions...${NC}"
chown -R "$USER:$GROUP" "$INSTALL_DIR"
chown -R "$USER:$GROUP" "$CONFIG_DIR"
chmod 600 "$CONFIG_DIR/config.yaml" 2>/dev/null || true
echo -e "${GREEN}  ✓ Permissions set${NC}"

# Install systemd service
echo -e "${YELLOW}Installing systemd service...${NC}"
cp systemd/whatsapp-proxy.service /etc/systemd/system/
systemctl daemon-reload
echo -e "${GREEN}  ✓ Service installed${NC}"

# Enable and start service
echo -e "${YELLOW}Enabling service...${NC}"
systemctl enable whatsapp-proxy
echo -e "${GREEN}  ✓ Service enabled (will start on boot)${NC}"

echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}Installation complete!${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo "Service commands:"
echo -e "  ${BLUE}Start:${NC}   sudo systemctl start whatsapp-proxy"
echo -e "  ${BLUE}Stop:${NC}    sudo systemctl stop whatsapp-proxy"
echo -e "  ${BLUE}Status:${NC}  sudo systemctl status whatsapp-proxy"
echo -e "  ${BLUE}Logs:${NC}    sudo journalctl -u whatsapp-proxy -f"
echo ""
echo -e "${YELLOW}Configuration:${NC} $CONFIG_DIR/config.yaml"
echo ""
echo -e "${YELLOW}Note:${NC} Service is enabled but not started."
echo -e "      Please edit the config file and start manually."
