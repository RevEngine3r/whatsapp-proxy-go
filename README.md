# WhatsApp Proxy Go

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)](https://github.com/RevEngine3r/whatsapp-proxy-go)

A lightweight, production-ready WhatsApp proxy server written in Go with single port operation and upstream SOCKS5 proxy support.

## ✨ Features

- 🚀 **Single Port Operation** - All protocols (HTTP, HTTPS, Jabber/XMPP) on one port
- 🔒 **SOCKS5 Upstream Support** - Route traffic through SOCKS5 proxy with authentication
- 🌍 **Cross-Platform** - Windows, Linux, macOS, FreeBSD
- 🏗️ **Multi-Architecture** - amd64, arm64, 386, arm
- ⚙️ **Flexible Configuration** - CLI arguments + YAML config + environment variables
- 📊 **Metrics & Monitoring** - OpenMetrics format (Prometheus-compatible)
- 🔐 **Auto SSL Certificates** - Self-signed certificate generation with caching
- 🎯 **Production Ready** - Systemd service, Windows service, Docker support
- ⚡ **High Performance** - Lightweight Go binary, efficient resource usage
- 🛡️ **Security Hardened** - Non-root execution, configurable limits

## 🎯 Key Differences from Original

This Go implementation reimagines the [original WhatsApp proxy](https://github.com/WhatsApp/proxy) with:

| Feature | Original (Docker) | This Implementation |
|---------|------------------|---------------------|
| **Ports** | Multiple (80, 443, 5222, 8080, 8443, etc.) | Single port (configurable) |
| **SOCKS5** | Not supported | Native SOCKS5 support |
| **Deployment** | Docker required | Native binary OR Docker |
| **Size** | ~200MB+ Docker image | ~10MB binary |
| **Platform** | Docker only | Windows, Linux, macOS, FreeBSD |
| **Architecture** | amd64 only | amd64, arm64, 386, arm |
| **Metrics** | Not available | OpenMetrics endpoint |
| **Service** | Docker/compose | Native systemd/Windows service |

## 📋 Requirements

**For Pre-built Binaries:**
- No dependencies required - single executable

**For Building from Source:**
- Go 1.21 or higher
- Make (optional)

**System Requirements:**
- CPU: 1 core (2+ recommended)
- RAM: 256MB minimum (512MB+ recommended)
- Disk: 50MB + logs
- Network: Stable internet connection

## 🚀 Quick Start

### Option 1: Download Pre-built Binary (Recommended)

```bash
# Linux amd64
wget https://github.com/RevEngine3r/whatsapp-proxy-go/releases/latest/download/whatsapp-proxy-linux-amd64
chmod +x whatsapp-proxy-linux-amd64
mv whatsapp-proxy-linux-amd64 whatsapp-proxy

# Create config
cp configs/config.example.yaml config.yaml
nano config.yaml  # Edit as needed

# Run
./whatsapp-proxy --config config.yaml
```

### Option 2: Using Docker

```bash
# Using docker-compose (recommended)
cp .env.example .env
cp configs/config.example.yaml config.yaml
docker-compose up -d

# Or using docker directly
docker run -d \
  --name whatsapp-proxy \
  -p 8443:8443 \
  -p 127.0.0.1:8199:8199 \
  -v $(pwd)/config.yaml:/etc/whatsapp-proxy/config.yaml:ro \
  whatsapp-proxy-go:latest
```

### Option 3: Build from Source

```bash
# Clone repository
git clone https://github.com/RevEngine3r/whatsapp-proxy-go.git
cd whatsapp-proxy-go

# Build
make build

# Or build for all platforms
./scripts/build-all.sh

# Run
./dist/whatsapp-proxy --config config.yaml
```

## 📖 Installation Guides

### Linux (systemd service)

```bash
# Install as service
sudo ./scripts/install-service-linux.sh

# Start service
sudo systemctl start whatsapp-proxy

# View logs
sudo journalctl -u whatsapp-proxy -f
```

### Windows (Windows service)

```powershell
# Run as Administrator
.\scripts\install-service-windows.bat

# Start service
sc start WhatsAppProxy

# Or use Services Manager (services.msc)
```

### macOS

```bash
# Simple run
./scripts/run.sh config.yaml

# Or install as launchd service (see INSTALL.md)
```

**📚 For detailed installation instructions, see [INSTALL.md](INSTALL.md)**

## ⚙️ Configuration

### Quick Configuration Examples

**Basic setup (no upstream proxy):**
```yaml
server:
  port: 8443
socks5:
  enabled: false
ssl:
  auto_generate: true
```

**With SOCKS5 upstream:**
```yaml
server:
  port: 8443
socks5:
  enabled: true
  host: 127.0.0.1
  port: 1080
  username: myuser
  password: mypass
ssl:
  auto_generate: true
logging:
  level: info
metrics:
  enabled: true
  port: 8199
```

**Production setup:**
```yaml
server:
  port: 443
  bind_addr: 0.0.0.0
  max_connections: 5000
  idle_timeout: 600
socks5:
  enabled: true
  host: proxy.example.com
  port: 1080
  username: ${SOCKS5_USER}
  password: ${SOCKS5_PASS}
ssl:
  auto_generate: false
  cert_file: /etc/ssl/certs/proxy.crt
  key_file: /etc/ssl/private/proxy.key
logging:
  level: warn
  format: json
  output: /var/log/whatsapp-proxy/proxy.log
metrics:
  enabled: true
  bind_addr: 127.0.0.1
```

### Configuration Methods

1. **YAML Config File** (recommended)
   ```bash
   ./whatsapp-proxy --config config.yaml
   ```

2. **Command-line Flags**
   ```bash
   ./whatsapp-proxy \
     --port 8443 \
     --socks5-proxy socks5://user:pass@127.0.0.1:1080 \
     --log-level info
   ```

3. **Environment Variables**
   ```bash
   export PROXY_PORT=8443
   export SOCKS5_ENABLED=true
   export SOCKS5_HOST=127.0.0.1
   ./whatsapp-proxy
   ```

**📚 For complete configuration reference, see [CONFIGURATION.md](CONFIGURATION.md)**

## 📊 Monitoring & Metrics

### Metrics Endpoint

Access metrics at `http://localhost:8199/metrics` (OpenMetrics/Prometheus format)

```bash
# View metrics
curl http://localhost:8199/metrics

# Health check
curl http://localhost:8199/health
```

### Available Metrics

- `whatsapp_proxy_connections_total` - Total connections
- `whatsapp_proxy_connections_active` - Active connections
- `whatsapp_proxy_connections_failed` - Failed connections
- `whatsapp_proxy_protocol_connections{protocol}` - Connections by protocol type
- `whatsapp_proxy_bytes_sent_total` - Total bytes sent
- `whatsapp_proxy_bytes_received_total` - Total bytes received
- `whatsapp_proxy_errors_total` - Total errors
- `whatsapp_proxy_uptime_seconds` - Server uptime

### Prometheus Configuration

```yaml
scrape_configs:
  - job_name: 'whatsapp-proxy'
    static_configs:
      - targets: ['localhost:8199']
```

## 🔧 Usage

### Running Directly

```bash
# Using config file
./whatsapp-proxy --config config.yaml

# Using CLI flags
./whatsapp-proxy --port 8443 --log-level debug

# Using helper scripts
./scripts/run.sh config.yaml    # Linux/macOS
.\scripts\run.bat config.yaml   # Windows
```

### Service Management

**Linux (systemd):**
```bash
sudo systemctl start whatsapp-proxy    # Start
sudo systemctl stop whatsapp-proxy     # Stop
sudo systemctl restart whatsapp-proxy  # Restart
sudo systemctl status whatsapp-proxy   # Status
sudo journalctl -u whatsapp-proxy -f   # Logs
```

**Windows:**
```powershell
sc start WhatsAppProxy     # Start
sc stop WhatsAppProxy      # Stop
sc query WhatsAppProxy     # Status
# Or use Services Manager (services.msc)
```

**Docker:**
```bash
docker-compose up -d         # Start
docker-compose down          # Stop
docker-compose restart       # Restart
docker-compose logs -f       # Logs
```

### Client Configuration

Configure your WhatsApp client to use the proxy:

**Android/iOS:**
1. Go to Settings → Storage and Data → Proxy
2. Enter proxy details:
   - Host: `your-server-ip`
   - Port: `8443` (or your configured port)
   - Type: HTTP or HTTPS

**Desktop:**
- Similar configuration in WhatsApp Desktop settings

## 🛠️ Development

### Build Commands

```bash
# Install dependencies
make deps

# Build for current platform
make build

# Build for all platforms
make build-all
# Or: ./scripts/build-all.sh

# Run tests
make test

# Run with coverage
make test-coverage

# Run linters
make lint

# Format code
make fmt

# Clean artifacts
make clean
```

### Project Structure

```
.
├── cmd/whatsapp-proxy/      # Main application
├── internal/
│   ├── config/              # Configuration management
│   ├── socks5/              # SOCKS5 client
│   ├── protocol/            # Protocol detection
│   ├── proxy/               # Proxy server core
│   └── ssl/                 # SSL certificate management
├── configs/                 # Configuration examples
├── scripts/                 # Helper scripts
├── systemd/                 # Systemd service files
├── ROAD_MAP/                # Development roadmap
└── docs/                    # Documentation
```

## 🔒 Security Considerations

### Best Practices

1. **Credentials**
   - Use environment variables for sensitive data
   - Never commit credentials to version control
   - Rotate passwords regularly

2. **Network Security**
   - Restrict metrics endpoint to localhost or trusted networks
   - Use firewall rules to limit access
   - Consider using VPN for proxy access

3. **SSL/TLS**
   - Use proper CA-signed certificates in production
   - Regularly rotate certificates
   - Self-signed certificates OK for internal/development use

4. **File Permissions**
   ```bash
   chmod 600 config.yaml              # Config file
   chmod 600 /etc/ssl/private/*.key   # Private keys
   chmod 755 /opt/whatsapp-proxy      # Binary directory
   ```

5. **System Hardening**
   - Run as non-root user (systemd service does this)
   - Set resource limits
   - Enable security features (SELinux, AppArmor)
   - Keep system and dependencies updated

## 📝 Project Status

✅ **Production Ready** - All core features implemented

### Completed Features

- ✅ Project setup and structure
- ✅ Configuration management (CLI + YAML + Env)
- ✅ SOCKS5 client implementation
- ✅ Proxy server core (HTTP/HTTPS/Jabber)
- ✅ SSL certificate management
- ✅ Deployment scripts and documentation

### Roadmap

- 📋 Docker Hub automated builds
- 📋 GitHub Actions CI/CD
- 📋 Automated release builds
- 📋 Package managers (Homebrew, apt, etc.)
- 📋 Web dashboard (optional)

See [PROGRESS.md](PROGRESS.md) for detailed development status.

## 🧪 Testing

### Run Tests

```bash
# All tests
make test

# With coverage
make test-coverage

# Specific package
go test ./internal/config/...

# Verbose
go test -v ./...

# Benchmarks
go test -bench=. ./internal/protocol/
```

### Test Coverage

- Config package: 11 tests
- SOCKS5 package: 11 tests (9 unit + 2 benchmarks)
- Protocol package: 8 tests (6 groups + 2 benchmarks)
- Proxy package: 4 tests
- SSL package: Multiple test scenarios
- **Total: 40+ tests** with high coverage

## 🤝 Contributing

Contributions are welcome! Here's how you can help:

1. **Report Bugs** - Open an issue with details
2. **Feature Requests** - Suggest new features
3. **Pull Requests** - Submit code improvements
4. **Documentation** - Improve docs and examples
5. **Testing** - Test on different platforms

### Development Workflow

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests (`make test`)
5. Commit your changes (`git commit -m 'Add amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Original [WhatsApp Proxy](https://github.com/WhatsApp/proxy) by Meta
- Go community for excellent libraries and tools
- Contributors and testers

## 📧 Support & Contact

- **Issues**: [GitHub Issues](https://github.com/RevEngine3r/whatsapp-proxy-go/issues)
- **Discussions**: [GitHub Discussions](https://github.com/RevEngine3r/whatsapp-proxy-go/discussions)
- **Documentation**: [INSTALL.md](INSTALL.md) | [CONFIGURATION.md](CONFIGURATION.md)
- **GitHub**: [@RevEngine3r](https://github.com/RevEngine3r)
- **Website**: [RevEngine3r.iR](https://www.RevEngine3r.iR)

## ⭐ Star History

If you find this project useful, please consider giving it a star! ⭐

---

**Disclaimer**: This is an unofficial reimplementation and is not affiliated with, endorsed by, or sponsored by WhatsApp or Meta. Use at your own risk and ensure compliance with local laws and WhatsApp's Terms of Service.
