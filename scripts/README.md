# Scripts Directory

Utility scripts for development, deployment, and service management.

## 🚀 Release Scripts

**Note**: Release scripts have moved to the repository root:
- [`github-release.bat`](../github-release.bat) - Create and push release tag
- [`build-release.ps1`](../build-release.ps1) - Build binaries locally
- [`build-release.bat`](../build-release.bat) - Windows launcher for build-release.ps1

See [RELEASE.md](../RELEASE.md) for release documentation.

## 💻 Development Scripts

### `run.bat` / `run.sh`
Quick run scripts for development.

**Windows**:
```batch
scripts\run.bat
```

**Linux/macOS**:
```bash
chmod +x scripts/run.sh
./scripts/run.sh
```

### `build-all.sh`
Build for multiple platforms (alternative to build-release.ps1).

```bash
chmod +x scripts/build-all.sh
./scripts/build-all.sh
```

## 🔧 Service Management

### Linux/macOS Service Installation

**Install**:
```bash
sudo ./scripts/install-service-linux.sh
```

**Uninstall**:
```bash
sudo ./scripts/uninstall-service-linux.sh
```

**Manage**:
```bash
sudo systemctl start whatsapp-proxy
sudo systemctl stop whatsapp-proxy
sudo systemctl status whatsapp-proxy
sudo systemctl enable whatsapp-proxy  # Start on boot
```

### Windows Service Installation

**Install** (run as Administrator):
```batch
scripts\install-service-windows.bat
```

**Manage**:
```batch
sc start WhatsAppProxy
sc stop WhatsAppProxy
sc query WhatsAppProxy
```

**Uninstall**:
```batch
sc delete WhatsAppProxy
```

## 📚 Documentation

- [Installation Guide](../INSTALL.md)
- [Configuration Guide](../CONFIGURATION.md)
- [Release Guide](../RELEASE.md)
- [Main README](../README.md)

## 🐛 Troubleshooting

### Script Permission Denied (Linux/macOS)
```bash
chmod +x scripts/*.sh
```

### PowerShell Execution Policy (Windows)
```powershell
Set-ExecutionPolicy RemoteSigned -Scope CurrentUser
```

### Service Won't Start
1. Check logs:
   - Linux: `journalctl -u whatsapp-proxy -f`
   - Windows: Event Viewer
2. Verify configuration file exists
3. Check port availability (default: 8080)

For more help, see the [main documentation](../README.md) or [open an issue](https://github.com/RevEngine3r/whatsapp-proxy-go/issues).
