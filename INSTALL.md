# Installation Guide

Comprehensive installation guide for WhatsApp Proxy Go across all supported platforms.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Quick Start](#quick-start)
- [Linux Installation](#linux-installation)
- [Windows Installation](#windows-installation)
- [macOS Installation](#macos-installation)
- [Docker Installation](#docker-installation)
- [Building from Source](#building-from-source)
- [Verification](#verification)
- [Troubleshooting](#troubleshooting)

## Prerequisites

### System Requirements

- **CPU**: 1 core minimum, 2+ cores recommended
- **RAM**: 256MB minimum, 512MB+ recommended
- **Disk**: 50MB for binary, additional space for logs
- **Network**: Stable internet connection
- **Ports**: One port for proxy (default 8443), one for metrics (default 8199)

### Software Requirements

- **Linux**: Any modern distribution with systemd (Ubuntu 18.04+, Debian 10+, CentOS 7+, etc.)
- **Windows**: Windows 10/11 or Windows Server 2016+
- **macOS**: macOS 10.15 (Catalina) or later
- **Go** (for building from source): Go 1.21 or later

## Quick Start

### Download Pre-built Binary

1. **Visit the releases page** (when available)
   ```bash
   # Example for Linux amd64
   wget https://github.com/RevEngine3r/whatsapp-proxy-go/releases/download/v1.0.0/whatsapp-proxy-linux-amd64
   chmod +x whatsapp-proxy-linux-amd64
   mv whatsapp-proxy-linux-amd64 whatsapp-proxy
   ```

2. **Create configuration file**
   ```bash
   # Copy example config
   cp configs/config.example.yaml config.yaml
   
   # Edit configuration
   nano config.yaml
   ```

3. **Run the proxy**
   ```bash
   ./whatsapp-proxy --config config.yaml
   ```

## Linux Installation

### Method 1: Systemd Service (Recommended)

Install as a system service that starts automatically on boot.

#### Step 1: Download or Build Binary

```bash
# Option A: Download release binary
wget https://github.com/RevEngine3r/whatsapp-proxy-go/releases/download/v1.0.0/whatsapp-proxy-linux-amd64
chmod +x whatsapp-proxy-linux-amd64
mv whatsapp-proxy-linux-amd64 whatsapp-proxy

# Option B: Build from source
git clone https://github.com/RevEngine3r/whatsapp-proxy-go.git
cd whatsapp-proxy-go
make build
```

#### Step 2: Configure

```bash
# Copy example configuration
cp configs/config.example.yaml config.yaml

# Edit configuration (see CONFIGURATION.md for details)
nano config.yaml
```

#### Step 3: Install Service

```bash
# Run installation script with sudo
sudo ./scripts/install-service-linux.sh
```

The script will:
- Create a dedicated `whatsapp-proxy` user
- Install binary to `/opt/whatsapp-proxy`
- Install config to `/etc/whatsapp-proxy`
- Install systemd service
- Enable service for auto-start

#### Step 4: Start Service

```bash
# Start the service
sudo systemctl start whatsapp-proxy

# Check status
sudo systemctl status whatsapp-proxy

# View logs
sudo journalctl -u whatsapp-proxy -f
```

#### Service Management

```bash
# Start service
sudo systemctl start whatsapp-proxy

# Stop service
sudo systemctl stop whatsapp-proxy

# Restart service
sudo systemctl restart whatsapp-proxy

# Enable auto-start on boot
sudo systemctl enable whatsapp-proxy

# Disable auto-start
sudo systemctl disable whatsapp-proxy

# View live logs
sudo journalctl -u whatsapp-proxy -f

# View logs since last boot
sudo journalctl -u whatsapp-proxy -b
```

### Method 2: Manual Run

Run directly without installing as a service.

```bash
# Using the run script
./scripts/run.sh config.yaml

# Or directly
./whatsapp-proxy --config config.yaml
```

### Uninstallation

```bash
# Run uninstall script
sudo ./scripts/uninstall-service-linux.sh
```

## Windows Installation

### Method 1: Windows Service (Recommended)

#### Step 1: Download or Build Binary

```powershell
# Download from releases (when available)
# Or build from source:
git clone https://github.com/RevEngine3r/whatsapp-proxy-go.git
cd whatsapp-proxy-go
make build
```

#### Step 2: Configure

```powershell
# Copy example configuration
copy configs\config.example.yaml config.yaml

# Edit with your preferred editor
notepad config.yaml
```

#### Step 3: Install Service

```powershell
# Run as Administrator
.\scripts\install-service-windows.bat
```

The script will:
- Create Windows service named "WhatsAppProxy"
- Configure auto-start on system boot
- Set up restart on failure

#### Step 4: Start Service

```powershell
# Using sc command
sc start WhatsAppProxy

# Or use Services Manager (services.msc)
# Find "WhatsApp Proxy Server" and click Start
```

#### Service Management

```powershell
# Start service
sc start WhatsAppProxy

# Stop service
sc stop WhatsAppProxy

# Query status
sc query WhatsAppProxy

# Remove service
sc delete WhatsAppProxy
```

**Using Services Manager GUI:**

1. Press `Win + R`, type `services.msc`, press Enter
2. Find "WhatsApp Proxy Server" in the list
3. Right-click for Start/Stop/Restart options
4. Double-click to configure startup type

### Method 2: Manual Run

Run directly without installing as a service.

```powershell
# Using the run script
.\scripts\run.bat config.yaml

# Or directly
.\whatsapp-proxy.exe --config config.yaml
```

## macOS Installation

### Using Homebrew (When Available)

```bash
# Future: Install via Homebrew
brew install whatsapp-proxy-go
```

### Manual Installation

#### Step 1: Download or Build

```bash
# Download release (when available)
curl -L -o whatsapp-proxy https://github.com/RevEngine3r/whatsapp-proxy-go/releases/download/v1.0.0/whatsapp-proxy-darwin-arm64
chmod +x whatsapp-proxy

# Or build from source
git clone https://github.com/RevEngine3r/whatsapp-proxy-go.git
cd whatsapp-proxy-go
make build
```

#### Step 2: Configure and Run

```bash
# Copy example config
cp configs/config.example.yaml config.yaml

# Edit configuration
nano config.yaml

# Run
./scripts/run.sh config.yaml
```

### Using launchd (Auto-start on Login)

Create `~/Library/LaunchAgents/com.whatsapp.proxy.plist`:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.whatsapp.proxy</string>
    <key>ProgramArguments</key>
    <array>
        <string>/usr/local/bin/whatsapp-proxy</string>
        <string>--config</string>
        <string>/usr/local/etc/whatsapp-proxy/config.yaml</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
    <key>StandardOutPath</key>
    <string>/usr/local/var/log/whatsapp-proxy.log</string>
    <key>StandardErrorPath</key>
    <string>/usr/local/var/log/whatsapp-proxy.error.log</string>
</dict>
</plist>
```

Load the service:

```bash
launchctl load ~/Library/LaunchAgents/com.whatsapp.proxy.plist
```

## Docker Installation

### Using Docker Compose (Recommended)

Create `docker-compose.yml`:

```yaml
version: '3.8'

services:
  whatsapp-proxy:
    image: whatsapp-proxy-go:latest
    container_name: whatsapp-proxy
    restart: unless-stopped
    ports:
      - "8443:8443"
      - "127.0.0.1:8199:8199"
    volumes:
      - ./config.yaml:/etc/whatsapp-proxy/config.yaml:ro
      - proxy-data:/data
    environment:
      - LOG_LEVEL=info
    networks:
      - proxy-network

volumes:
  proxy-data:

networks:
  proxy-network:
    driver: bridge
```

Run:

```bash
docker-compose up -d
```

### Using Docker CLI

```bash
docker run -d \
  --name whatsapp-proxy \
  --restart unless-stopped \
  -p 8443:8443 \
  -p 127.0.0.1:8199:8199 \
  -v $(pwd)/config.yaml:/etc/whatsapp-proxy/config.yaml:ro \
  whatsapp-proxy-go:latest
```

## Building from Source

### Prerequisites

- Go 1.21 or later
- Git
- Make (optional, but recommended)

### Build Steps

```bash
# Clone repository
git clone https://github.com/RevEngine3r/whatsapp-proxy-go.git
cd whatsapp-proxy-go

# Install dependencies
go mod download

# Build for current platform
make build

# Or build manually
go build -o dist/whatsapp-proxy ./cmd/whatsapp-proxy

# Build for all platforms
./scripts/build-all.sh
```

### Build Targets

```bash
# Build for current platform
make build

# Build for Linux amd64
make build-linux

# Build for Windows
make build-windows

# Build for macOS
make build-darwin

# Build for all platforms
make build-all

# Run tests
make test

# Clean build artifacts
make clean
```

## Verification

### Check if Proxy is Running

```bash
# Check if port is listening
netstat -tuln | grep 8443

# Test health endpoint
curl http://localhost:8199/health
# Expected output: OK

# Check metrics
curl http://localhost:8199/metrics
```

### Test Proxy Connection

```bash
# Test HTTP proxy
curl -x http://localhost:8443 http://example.com

# Test HTTPS proxy
curl -x http://localhost:8443 https://example.com
```

### View Logs

```bash
# Linux (systemd)
sudo journalctl -u whatsapp-proxy -f

# Windows (Event Viewer)
# Open Event Viewer > Application logs

# Manual run
# Check console output or configured log file
```

## Troubleshooting

### Proxy Won't Start

1. **Check if port is already in use**
   ```bash
   # Linux/macOS
   netstat -tuln | grep 8443
   lsof -i :8443
   
   # Windows
   netstat -ano | findstr :8443
   ```

2. **Check configuration file**
   ```bash
   # Validate YAML syntax
   ./whatsapp-proxy --config config.yaml --validate
   ```

3. **Check permissions**
   ```bash
   # Linux: Ensure binary is executable
   chmod +x whatsapp-proxy
   
   # Ensure config file is readable
   chmod 644 config.yaml
   ```

### Connection Issues

1. **Check firewall rules**
   ```bash
   # Linux (Ubuntu/Debian)
   sudo ufw status
   sudo ufw allow 8443/tcp
   
   # Linux (CentOS/RHEL)
   sudo firewall-cmd --list-all
   sudo firewall-cmd --permanent --add-port=8443/tcp
   sudo firewall-cmd --reload
   
   # Windows
   # Open Windows Firewall > Advanced Settings
   # Add inbound rule for port 8443
   ```

2. **Test SOCKS5 upstream connectivity**
   ```bash
   # Try connecting to SOCKS5 proxy directly
   curl -x socks5://localhost:1080 https://example.com
   ```

3. **Check logs for errors**
   ```bash
   # Increase log level to debug
   # Edit config.yaml: logging.level = debug
   ```

### Performance Issues

1. **Increase connection limits**
   ```yaml
   # In config.yaml
   server:
     max_connections: 5000
   ```

2. **Check system limits**
   ```bash
   # Linux: Check file descriptor limits
   ulimit -n
   
   # Increase if needed (add to /etc/security/limits.conf)
   * soft nofile 65535
   * hard nofile 65535
   ```

3. **Monitor metrics**
   ```bash
   # Watch metrics in real-time
   watch -n 1 'curl -s http://localhost:8199/metrics | grep connections'
   ```

### Certificate Issues

1. **Regenerate certificates**
   ```bash
   # Remove cached certificates
   rm -rf ~/.cache/whatsapp-proxy/certs/
   
   # Restart proxy to regenerate
   ```

2. **Use custom certificates**
   ```yaml
   # In config.yaml
   ssl:
     auto_generate: false
     cert_file: /path/to/cert.pem
     key_file: /path/to/key.pem
   ```

### Getting Help

- **Documentation**: [README.md](README.md), [CONFIGURATION.md](CONFIGURATION.md)
- **Issues**: [GitHub Issues](https://github.com/RevEngine3r/whatsapp-proxy-go/issues)
- **Logs**: Enable debug logging for detailed troubleshooting

## Next Steps

- Review [CONFIGURATION.md](CONFIGURATION.md) for detailed configuration options
- Set up monitoring using the metrics endpoint
- Configure automatic backups of your configuration
- Set up log rotation for production deployments
