# Release Workflow

This document describes the automated release process for WhatsApp Proxy Go.

## Overview

The project uses GitHub Actions for automated releases with datetime-based versioning. Releases are triggered by pushing a special `release` tag.

## Version Format

Versions follow the format: `YYYY.MM.DD.HHMM`

**Example**: `2026.01.08.1530` (January 8, 2026, 15:30 UTC)

## Supported Platforms

The CI/CD builds binaries for the following platforms:

### Linux
- `linux-amd64` - 64-bit Intel/AMD
- `linux-arm64` - 64-bit ARM (Raspberry Pi 4, AWS Graviton)
- `linux-386` - 32-bit Intel/AMD
- `linux-armv7` - 32-bit ARM (Raspberry Pi 2/3)

### macOS
- `darwin-amd64` - Intel-based Macs
- `darwin-arm64` - Apple Silicon (M1/M2/M3)

### Windows
- `windows-amd64` - 64-bit Intel/AMD
- `windows-386` - 32-bit Intel/AMD
- `windows-arm64` - ARM64 (Surface Pro X)

## Creating a Release

### Prerequisites

1. Ensure all changes are committed and pushed to `main`
2. All tests are passing
3. You have push access to the repository

### Release Steps

#### 1. Create and Push the Release Tag

```bash
# Make sure you're on the main branch and up to date
git checkout main
git pull origin main

# Create the release tag
git tag release

# Push the tag to trigger the workflow
git push origin release
```

#### 2. Monitor the Workflow

1. Go to the [Actions tab](https://github.com/RevEngine3r/whatsapp-proxy-go/actions)
2. Watch the "Release" workflow progress
3. The workflow will:
   - Generate a datetime-based version
   - Run all tests
   - Build binaries for all platforms
   - Generate SHA256 checksums
   - Create a GitHub release

#### 3. Verify the Release

1. Check the [Releases page](https://github.com/RevEngine3r/whatsapp-proxy-go/releases)
2. Verify all binary files are present
3. Download and test the binary for your platform

### Cleanup After Release

After a successful release, delete the local and remote `release` tag:

```bash
# Delete local tag
git tag -d release

# Delete remote tag
git push origin :refs/tags/release
```

This allows you to reuse the same tag for future releases.

## Build Information

Each binary is built with the following information embedded:

- **Version**: Datetime-based version number
- **Build Time**: UTC timestamp of the build
- **Git Commit**: Short SHA of the commit

These can be accessed via version flags in the application.

## Binary Naming Convention

Binaries follow this naming pattern:

```
whatsapp-proxy-{VERSION}-{OS}-{ARCH}[v{ARM_VERSION}].{EXT}
```

**Examples**:
- `whatsapp-proxy-2026.01.08.1530-linux-amd64.tar.gz`
- `whatsapp-proxy-2026.01.08.1530-darwin-arm64.tar.gz`
- `whatsapp-proxy-2026.01.08.1530-windows-amd64.zip`
- `whatsapp-proxy-2026.01.08.1530-linux-armv7.tar.gz`

## Checksum Verification

Every release includes a `checksums.txt` file with SHA256 hashes.

### Verify on Linux/macOS

```bash
# Download the checksum file and your binary
wget https://github.com/RevEngine3r/whatsapp-proxy-go/releases/download/v2026.01.08.1530/checksums.txt
wget https://github.com/RevEngine3r/whatsapp-proxy-go/releases/download/v2026.01.08.1530/whatsapp-proxy-2026.01.08.1530-linux-amd64.tar.gz

# Verify (will show "OK" if valid)
sha256sum -c checksums.txt --ignore-missing
```

### Verify on Windows (PowerShell)

```powershell
# Calculate hash of downloaded file
$hash = (Get-FileHash whatsapp-proxy-2026.01.08.1530-windows-amd64.zip -Algorithm SHA256).Hash

# Compare with checksums.txt
Get-Content checksums.txt | Select-String -Pattern $hash
```

## Troubleshooting

### Workflow Fails

1. Check the [Actions tab](https://github.com/RevEngine3r/whatsapp-proxy-go/actions) for error details
2. Common issues:
   - Test failures: Fix tests and retry
   - Build errors: Check Go version compatibility
   - Permission issues: Verify repository settings

### Tag Already Exists

If the `release` tag wasn't cleaned up:

```bash
# Force delete and recreate
git tag -d release
git push origin :refs/tags/release
git tag release
git push origin release
```

### Missing Binaries

If some platform binaries are missing:

1. Check the workflow logs for build errors
2. Verify the platform is listed in the workflow
3. Re-run the failed jobs from the Actions tab

## Manual Release (Alternative)

If you need to create a release manually:

```bash
# Set version
VERSION=$(date -u +'%Y.%m.%d.%H%M')

# Build for specific platform
GOOS=linux GOARCH=amd64 go build \
  -ldflags="-s -w -X main.version=${VERSION}" \
  -o whatsapp-proxy-${VERSION}-linux-amd64 \
  ./cmd/whatsapp-proxy

# Create archive
tar -czf whatsapp-proxy-${VERSION}-linux-amd64.tar.gz whatsapp-proxy-${VERSION}-linux-amd64

# Generate checksum
sha256sum whatsapp-proxy-${VERSION}-linux-amd64.tar.gz > checksum.txt
```

## CI/CD Configuration

The workflow is defined in `.github/workflows/release.yml`.

Key features:
- **Trigger**: Push of `release` tag
- **Go Version**: 1.21+ (defined in workflow)
- **Tests**: Full test suite with race detection
- **Coverage**: Uploaded to Codecov (optional)
- **Build Flags**: `-trimpath`, `-ldflags="-s -w"` for smaller binaries

## Best Practices

1. **Test Before Release**: Always test locally before pushing the release tag
2. **Review Changes**: Check git log to ensure all intended changes are included
3. **Update Docs**: Update README/CHANGELOG before releasing
4. **Version Tags**: Keep release tags clean by deleting after use
5. **Communication**: Announce releases in appropriate channels

## Additional Resources

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Go Cross-Compilation](https://go.dev/doc/install/source#environment)
- [Semantic Versioning](https://semver.org/) (for reference, though we use datetime)

## Questions?

For issues with the release process:
1. Check this documentation first
2. Review closed issues in the repository
3. Open a new issue with the `ci/cd` label
