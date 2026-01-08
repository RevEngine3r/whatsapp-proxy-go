#!/bin/bash
# WhatsApp Proxy Go - Linux Service Uninstallation Script

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Check if running as root
if [ "$EUID" -ne 0 ]; then
    echo -e "${RED}Error: This script must be run as root${NC}"
    echo "Usage: sudo $0"
    exit 1
fi

SERVICE_NAME="whatsapp-proxy"
USER="whatsapp-proxy"
INSTALL_DIR="/opt/whatsapp-proxy"
CONFIG_DIR="/etc/whatsapp-proxy"

echo -e "${YELLOW}Uninstalling WhatsApp Proxy service...${NC}"
echo ""

# Stop and disable service
if systemctl is-active --quiet "$SERVICE_NAME"; then
    echo -e "${YELLOW}Stopping service...${NC}"
    systemctl stop "$SERVICE_NAME"
    echo -e "${GREEN}  ✓ Service stopped${NC}"
fi

if systemctl is-enabled --quiet "$SERVICE_NAME" 2>/dev/null; then
    echo -e "${YELLOW}Disabling service...${NC}"
    systemctl disable "$SERVICE_NAME"
    echo -e "${GREEN}  ✓ Service disabled${NC}"
fi

# Remove service file
if [ -f "/etc/systemd/system/$SERVICE_NAME.service" ]; then
    echo -e "${YELLOW}Removing service file...${NC}"
    rm /etc/systemd/system/$SERVICE_NAME.service
    systemctl daemon-reload
    echo -e "${GREEN}  ✓ Service file removed${NC}"
fi

# Ask before removing files
read -p "Remove installation directory $INSTALL_DIR? [y/N] " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    rm -rf "$INSTALL_DIR"
    echo -e "${GREEN}  ✓ Installation directory removed${NC}"
fi

read -p "Remove configuration directory $CONFIG_DIR? [y/N] " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    rm -rf "$CONFIG_DIR"
    echo -e "${GREEN}  ✓ Configuration directory removed${NC}"
fi

read -p "Remove service user $USER? [y/N] " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    userdel "$USER" 2>/dev/null || true
    echo -e "${GREEN}  ✓ Service user removed${NC}"
fi

echo ""
echo -e "${GREEN}Uninstallation complete!${NC}"
