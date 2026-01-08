# WhatsApp Proxy Go - Project Map

## Overview
A lightweight, cross-platform WhatsApp proxy server written in Go that supports:
- Single port operation (simplifies deployment)
- Upstream SOCKS5 proxy support
- Multi-protocol handling (HTTP, HTTPS, Jabber/XMPP)
- Cross-platform binaries (Windows, Linux, macOS, BSD)
- Multiple architectures (amd64, arm64, 386, arm)

## Architecture

### Core Components
1. **Config Manager** - Handles configuration from CLI args and config file
2. **SOCKS5 Client** - Upstream proxy connection management
3. **Proxy Server** - Main server handling all protocols on single port
4. **SSL Manager** - Certificate generation and management
5. **Protocol Router** - Routes traffic based on protocol detection

### Key Features
- Protocol detection (HTTP/HTTPS/Jabber)
- Dynamic certificate generation
- Graceful shutdown
- Metrics endpoint
- Logging with levels

## Directory Structure
```
whatsapp-proxy-go/
├── cmd/
│   └── whatsapp-proxy/
│       └── main.go
├── internal/
│   ├── config/
│   ├── proxy/
│   ├── socks5/
│   ├── ssl/
│   └── protocol/
├── configs/
│   └── config.example.yaml
├── scripts/
│   ├── run.sh
│   ├── run.bat
│   ├── install-service-linux.sh
│   └── install-service-windows.bat
├── systemd/
│   └── whatsapp-proxy.service
├── deployments/
│   └── docker/
├── ROAD_MAP/
├── go.mod
└── README.md
```

## Technology Stack
- **Language**: Go 1.21+
- **Key Libraries**:
  - Standard library (net/http, crypto/tls)
  - golang.org/x/net for advanced networking
  - gopkg.in/yaml.v3 for config parsing
  - github.com/spf13/cobra for CLI
  - github.com/spf13/viper for config management

## Development Phases
See ROAD_MAP/README.md for detailed feature breakdown and implementation steps.
