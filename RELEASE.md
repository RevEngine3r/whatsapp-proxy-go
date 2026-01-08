# Release Guide

Simple guide for creating releases.

## Quick Start

### Windows

1. **Run the release script**:
   ```batch
   github-release.bat
   ```

2. **Monitor the build**:
   - GitHub Actions will automatically build binaries
   - Check progress at: https://github.com/RevEngine3r/whatsapp-proxy-go/actions

3. **Download release**:
   - Visit: https://github.com/RevEngine3r/whatsapp-proxy-go/releases

### Manual (Any OS)

```bash
# Create and push release tag
git tag release
git push origin release --force

# After workflow completes (optional cleanup)
git tag -d release
git push origin :refs/tags/release
```

## Local Build

To build locally without creating a release:

```bash
# Windows
build-release.bat

# Linux/macOS
pwsh ./build-release.ps1
```

Builds will be created in the `dist/` directory.

## What Happens

1. **Version Generation**: Automatic datetime version (YYYY.MM.DD.HHmm)
2. **Multi-Platform Build**: Binaries for Linux, macOS, Windows
3. **Architectures**: amd64, arm64, 386, armv7
4. **Checksums**: SHA256 for all files
5. **GitHub Release**: Automatic creation with all files

## Files Generated

- `whatsapp-proxy-VERSION-OS-ARCH.tar.gz` (Linux/macOS)
- `whatsapp-proxy-VERSION-windows-ARCH.zip` (Windows)
- `checksums.txt` (SHA256 hashes)

## Troubleshooting

| Problem | Solution |
|---------|----------|
| Git not found | Install Git and add to PATH |
| PowerShell error | Install PowerShell 7+ |
| Build fails | Check Go is installed (1.21+) |
| Tests fail | Fix tests before releasing |

For more details, see [README.md](README.md) and [INSTALL.md](INSTALL.md).
