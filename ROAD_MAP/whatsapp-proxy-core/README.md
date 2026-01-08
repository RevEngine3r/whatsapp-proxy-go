# Feature: WhatsApp Proxy Core Implementation

## Overview
Implement a complete WhatsApp proxy server in Go that operates on a single port and supports upstream SOCKS5 proxy connections. The implementation must be cross-platform and production-ready.

## Goals
1. Single port operation handling all protocols
2. SOCKS5 upstream proxy support
3. Cross-platform compatibility (Windows, Linux, macOS, BSD)
4. Multi-architecture support (amd64, arm64, 386, arm)
5. Flexible configuration (CLI args + config file)
6. Production-ready deployment scripts

## Steps

### Step 1: Project Setup and Structure
- Initialize Go module
- Create project directory structure
- Set up build system with cross-compilation
- Add essential dependencies
- Create basic README

### Step 2: Configuration Management
- Implement config structure
- Add CLI argument parsing with cobra
- Add YAML config file support with viper
- Configuration validation
- Environment variable support

### Step 3: SOCKS5 Client Implementation
- SOCKS5 protocol handler
- Connection pooling
- Authentication support
- Error handling and retries
- Connection testing

### Step 4: Proxy Server Core
- Protocol detection (HTTP/HTTPS/Jabber)
- Single port listener
- Request routing
- Connection handling
- Upstream forwarding via SOCKS5
- Metrics endpoint

### Step 5: SSL Certificate Management
- Self-signed certificate generation
- Certificate caching
- TLS configuration
- Certificate rotation support

### Step 6: Deployment and Service Scripts
- Shell script runner (Linux/macOS)
- Batch script runner (Windows)
- Systemd service file
- Windows service installer
- Example configuration file
- Docker support (optional)

## Success Criteria
- [ ] Proxy handles HTTP, HTTPS, and Jabber on single port
- [ ] Successfully connects via upstream SOCKS5 proxy
- [ ] Cross-compiles for all major OS/arch combinations
- [ ] Configuration works via both CLI and config file
- [ ] Service scripts install and run correctly
- [ ] Comprehensive tests pass
- [ ] Documentation complete
