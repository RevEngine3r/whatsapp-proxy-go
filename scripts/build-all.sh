#!/bin/bash
# WhatsApp Proxy Go - Cross-Platform Build Script
# Builds binaries for all supported platforms

set -e

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Configuration
APP_NAME="whatsapp-proxy"
OUT_DIR="dist"
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS="-s -w -X main.Version=$VERSION -X main.BuildTime=$BUILD_TIME"

# Platforms to build
PLATFORMS=(
    "linux/amd64"
    "linux/arm64"
    "linux/386"
    "linux/arm"
    "windows/amd64"
    "windows/386"
    "darwin/amd64"
    "darwin/arm64"
    "freebsd/amd64"
)

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  WhatsApp Proxy Go - Build Script${NC}"
echo -e "${BLUE}========================================${NC}"
echo -e "${YELLOW}Version:${NC} $VERSION"
echo -e "${YELLOW}Build Time:${NC} $BUILD_TIME"
echo ""

# Clean and create output directory
rm -rf "$OUT_DIR"
mkdir -p "$OUT_DIR"

echo -e "${GREEN}Building for ${#PLATFORMS[@]} platforms...${NC}"
echo ""

# Build for each platform
for platform in "${PLATFORMS[@]}"; do
    GOOS=${platform%/*}
    GOARCH=${platform#*/}
    
    output="$OUT_DIR/${APP_NAME}-${GOOS}-${GOARCH}"
    
    # Add .exe extension for Windows
    if [ "$GOOS" = "windows" ]; then
        output="${output}.exe"
    fi
    
    echo -e "${YELLOW}Building:${NC} $GOOS/$GOARCH"
    
    # Build
    GOOS=$GOOS GOARCH=$GOARCH go build \
        -ldflags="$LDFLAGS" \
        -o "$output" \
        ./cmd/whatsapp-proxy
    
    if [ $? -eq 0 ]; then
        size=$(du -h "$output" | cut -f1)
        echo -e "${GREEN}  ✓ Built:${NC} $output (${size})"
    else
        echo -e "${RED}  ✗ Failed:${NC} $output"
        exit 1
    fi
    echo ""
done

# Generate checksums
echo -e "${YELLOW}Generating checksums...${NC}"
cd "$OUT_DIR"
sha256sum * > SHA256SUMS
cd ..

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}Build complete!${NC}"
echo -e "${GREEN}========================================${NC}"
echo -e "Output directory: ${BLUE}$OUT_DIR${NC}"
echo ""
echo "Binaries:"
ls -lh "$OUT_DIR" | grep -v "^total" | grep -v "SHA256SUMS"
echo ""
echo -e "${YELLOW}Checksums saved to:${NC} $OUT_DIR/SHA256SUMS"
