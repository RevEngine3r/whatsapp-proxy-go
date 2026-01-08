# WhatsApp Proxy Go

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)](https://github.com/RevEngine3r/whatsapp-proxy-go)

A lightweight, cross-platform WhatsApp proxy server written in Go with single port operation and upstream SOCKS5 proxy support.

## ✨ Features

- 🚀 **Single Port Operation** - All protocols (HTTP, HTTPS, Jabber/XMPP) on one port
- 🔒 **SOCKS5 Upstream Support** - Route traffic through SOCKS5 proxy with authentication
- 🌍 **Cross-Platform** - Windows, Linux, macOS, FreeBSD
- 🏗️ **Multi-Architecture** - amd64, arm64, 386, arm
- ⚙️ **Flexible Configuration** - CLI arguments + YAML config file
- 📊 **Metrics Endpoint** - OpenMetrics format for monitoring
- 🔐 **Auto SSL Certificates** - Self-signed certificate generation
- 🎯 **Production Ready** - Systemd service, Windows service support

## 🎯 Key Differences from Original

This Go implementation reimagines the [original WhatsApp proxy](https://github.com/WhatsApp/proxy) with:

1. **Single Port** instead of multiple ports (80, 443, 5222, 8080, 8443, etc.)
2. **Native SOCKS5 Support** for upstream proxy routing
3. **Go Binary** - No Docker required, single executable
4. **Smaller Footprint** - Lightweight native binary

## 📋 Requirements

- Go 1.21 or higher (for building from source)
- OR download pre-built binaries from releases

## 🚀 Quick Start

### Using Pre-built Binaries

1. Download the latest release for your platform:
```bash
# Linux amd64
wget https://github.com/RevEngine3r/whatsapp-proxy-go/releases/latest/download/whatsapp-proxy-linux-amd64
chmod +x whatsapp-proxy-linux-amd64

# Windows amd64
# Download whatsapp-proxy-windows-amd64.exe from releases
```

2. Create a configuration file:
```bash
cp configs/config.example.yaml config.yaml
# Edit config.yaml with your settings
```

3. Run the proxy:
```bash
# Linux/macOS
./whatsapp-proxy-linux-amd64 --config config.yaml

# Windows
whatsapp-proxy-windows-amd64.exe --config config.yaml
```

### Building from Source

```bash
# Clone repository
git clone https://github.com/RevEngine3r/whatsapp-proxy-go.git
cd whatsapp-proxy-go

# Build for current platform
make build

# Build for all platforms
make build-all

# Run
./dist/whatsapp-proxy --config config.yaml
```

## ⚙️ Configuration

### CLI Arguments

```bash
whatsapp-proxy --port 8443 \
               --bind 0.0.0.0 \
               --socks5-proxy socks5://user:pass@127.0.0.1:1080 \
               --log-level info
```

### Configuration File

See `configs/config.example.yaml` for a complete example with all options.

```yaml
server:
  port: 8443
  bind_addr: 0.0.0.0

socks5:
  enabled: true
  host: 127.0.0.1
  port: 1080
  username: ""
  password: ""

ssl:
  auto_generate: true

logging:
  level: info

metrics:
  enabled: true
  port: 8199
```

## 🔧 Installation as Service

### Linux (systemd)

```bash
# Copy files
sudo mkdir -p /opt/whatsapp-proxy /etc/whatsapp-proxy
sudo cp dist/whatsapp-proxy /opt/whatsapp-proxy/
sudo cp configs/config.example.yaml /etc/whatsapp-proxy/config.yaml

# Install service
sudo cp systemd/whatsapp-proxy.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable whatsapp-proxy
sudo systemctl start whatsapp-proxy
```

### Windows Service

```batch
REM Run as Administrator
sc create WhatsAppProxy binPath= "%CD%\whatsapp-proxy.exe --config %CD%\config.yaml" start= auto
sc start WhatsAppProxy
```

## 📊 Monitoring

Access metrics at `http://<host>:8199/metrics` (OpenMetrics format)

```bash
curl http://localhost:8199/metrics
```

## 🔒 Security Considerations

- Always use strong passwords for SOCKS5 authentication
- Restrict metrics endpoint to localhost or trusted networks
- Use custom SSL certificates in production
- Keep the binary updated
- Monitor logs for suspicious activity

## 📚 Documentation

Detailed documentation is available in the `docs/` directory:

- [Installation Guide](docs/INSTALL.md) *(coming soon)*
- [Configuration Reference](docs/CONFIGURATION.md) *(coming soon)*
- [Troubleshooting](docs/TROUBLESHOOTING.md) *(coming soon)*

## 🛠️ Development

```bash
# Install dependencies
make deps

# Run tests
make test

# Run linters
make lint

# Format code
make fmt

# Build
make build
```

## 📝 Project Status

🚧 **In Development** - Step 1 of 6 completed

See [PROGRESS.md](PROGRESS.md) for current implementation status.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Original [WhatsApp Proxy](https://github.com/WhatsApp/proxy) by Meta
- Inspired by the need for a simpler, Go-native implementation

## 📧 Contact

- GitHub: [@RevEngine3r](https://github.com/RevEngine3r)
- Website: [RevEngine3r.iR](https://www.RevEngine3r.iR)

---

**Note**: This is an unofficial reimplementation and is not affiliated with, endorsed by, or sponsored by WhatsApp or Meta.
