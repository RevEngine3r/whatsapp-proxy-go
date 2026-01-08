# Step 6: Deployment and Service Scripts

## Objective
Create production-ready deployment scripts for all platforms, service installation, and comprehensive documentation.

## Tasks

### 1. Runner Scripts
- Shell script for Linux/macOS (run.sh)
- Batch script for Windows (run.bat)
- Configuration file checking
- Basic error handling

### 2. Linux Service Integration
- Systemd service file
- Installation script
- Enable on boot
- Log management

### 3. Windows Service
- Service installer batch
- NSSM or sc.exe integration
- Registry configuration
- Event log integration

### 4. Example Configuration
- Complete config.example.yaml
- Inline documentation
- Common scenarios
- Security notes

### 5. Build and Release
- Cross-compilation script
- Release packaging
- Checksum generation
- Version tagging

### 6. Documentation
- Complete README
- Installation guide
- Configuration reference
- Troubleshooting
- Security considerations

## Implementation Details

### run.sh (Linux/macOS)
```bash
#!/bin/bash
set -e

CONFIG_FILE="${CONFIG_FILE:-config.yaml}"
LOG_FILE="${LOG_FILE:-whatsapp-proxy.log}"

if [ ! -f "$CONFIG_FILE" ]; then
    echo "Config file not found: $CONFIG_FILE"
    exit 1
fi

./whatsapp-proxy --config "$CONFIG_FILE" 2>&1 | tee -a "$LOG_FILE"
```

### run.bat (Windows)
```batch
@echo off
setlocal

set CONFIG_FILE=config.yaml
if not exist "%CONFIG_FILE%" (
    echo Config file not found: %CONFIG_FILE%
    exit /b 1
)

whatsapp-proxy.exe --config %CONFIG_FILE%
```

### Systemd Service (whatsapp-proxy.service)
```ini
[Unit]
Description=WhatsApp Proxy Server
After=network.target

[Service]
Type=simple
User=whatsapp-proxy
Group=whatsapp-proxy
WorkingDirectory=/opt/whatsapp-proxy
ExecStart=/opt/whatsapp-proxy/whatsapp-proxy --config /etc/whatsapp-proxy/config.yaml
Restart=on-failure
RestartSec=5
StandardOutput=journal
StandardError=journal
SyslogIdentifier=whatsapp-proxy

[Install]
WantedBy=multi-user.target
```

### Linux Install Script
```bash
#!/bin/bash
# install-service-linux.sh

sudo useradd -r -s /bin/false whatsapp-proxy || true
sudo mkdir -p /opt/whatsapp-proxy /etc/whatsapp-proxy
sudo cp whatsapp-proxy /opt/whatsapp-proxy/
sudo cp config.example.yaml /etc/whatsapp-proxy/config.yaml
sudo cp systemd/whatsapp-proxy.service /etc/systemd/system/
sudo chown -R whatsapp-proxy:whatsapp-proxy /opt/whatsapp-proxy
sudo systemctl daemon-reload
sudo systemctl enable whatsapp-proxy
sudo systemctl start whatsapp-proxy
```

### Windows Service Install
```batch
REM install-service-windows.bat
sc create WhatsAppProxy binPath= "%CD%\whatsapp-proxy.exe --config %CD%\config.yaml" start= auto
sc description WhatsAppProxy "WhatsApp Proxy Server"
sc start WhatsAppProxy
```

### Build Script
```bash
#!/bin/bash
# build-all.sh

VERSION=$(git describe --tags --always)
LDFLAGS="-s -w -X main.Version=$VERSION"

platforms=(
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

for platform in "${platforms[@]}"; do
    GOOS=${platform%/*}
    GOARCH=${platform#*/}
    output="whatsapp-proxy-${GOOS}-${GOARCH}"
    [ "$GOOS" = "windows" ] && output="${output}.exe"
    
    GOOS=$GOOS GOARCH=$GOARCH go build -ldflags="$LDFLAGS" -o "dist/$output" ./cmd/whatsapp-proxy
    echo "Built: $output"
done
```

### Example Config (Complete)
```yaml
# WhatsApp Proxy Configuration

server:
  # Port to listen on (single port for all protocols)
  port: 8443
  
  # Bind address (0.0.0.0 for all interfaces)
  bind_addr: 0.0.0.0
  
  # Idle connection timeout (seconds)
  idle_timeout: 300
  
  # Maximum concurrent connections
  max_connections: 1000

# Upstream SOCKS5 proxy configuration
socks5:
  # Enable upstream SOCKS5 proxy
  enabled: true
  
  # SOCKS5 proxy host
  host: 127.0.0.1
  
  # SOCKS5 proxy port
  port: 1080
  
  # Username (leave empty for no auth)
  username: ""
  
  # Password (leave empty for no auth)
  password: ""
  
  # Connection timeout (seconds)
  timeout: 30

# SSL/TLS configuration
ssl:
  # Auto-generate self-signed certificates
  auto_generate: true
  
  # Custom certificate file (if not auto-generating)
  cert_file: ""
  
  # Custom key file (if not auto-generating)
  key_file: ""
  
  # DNS names for certificate SANs
  dns_names:
    - localhost
    - proxy.example.com
  
  # IP addresses for certificate SANs
  ip_addresses:
    - 127.0.0.1
  
  # Certificate validity in days
  validity_days: 365

# Logging configuration
logging:
  # Log level: debug, info, warn, error
  level: info
  
  # Log format: text, json
  format: text
  
  # Log output: stdout, stderr, or file path
  output: stdout

# Metrics endpoint
metrics:
  # Enable metrics endpoint
  enabled: true
  
  # Metrics server port
  port: 8199
  
  # Bind address for metrics
  bind_addr: 127.0.0.1
```

## Testing
- Scripts run without errors
- Service installs correctly
- Service starts on boot
- Config file loads properly
- Cross-platform builds succeed
- Documentation accurate

## Deliverables
- scripts/run.sh
- scripts/run.bat
- scripts/install-service-linux.sh
- scripts/install-service-windows.bat
- scripts/build-all.sh
- systemd/whatsapp-proxy.service
- configs/config.example.yaml
- Complete README.md
- INSTALL.md
- CONFIGURATION.md
