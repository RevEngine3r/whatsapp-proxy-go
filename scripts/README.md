# Scripts Directory

This directory contains utility scripts for automating common development and release tasks.

## 📁 Available Scripts

### Release Automation

#### `release.bat` - Windows Batch Release Script
**Purpose**: Automates the creation and pushing of release tags for Windows users.

**Features**:
- ✅ Validates Git installation
- ✅ Checks repository status
- ✅ Fetches latest changes
- ✅ Warns about uncommitted changes
- ✅ Displays release summary with version preview
- ✅ Cleans up old release tags automatically
- ✅ Pushes release tag to trigger CI/CD
- ✅ Provides helpful next-step instructions

**Usage**:
```batch
cd /path/to/whatsapp-proxy-go
scripts\release.bat
```

**Requirements**:
- Windows (any version)
- Git for Windows installed
- Git executable in PATH

---

#### `release.ps1` - PowerShell Release Script
**Purpose**: Modern PowerShell version with enhanced features and better error handling.

**Features**:
- ✅ All features from `release.bat`
- ✅ Colored console output
- ✅ Better error handling and reporting
- ✅ Option to open GitHub Actions in browser
- ✅ More detailed status information
- ✅ PowerShell 5.1+ compatible

**Usage**:
```powershell
cd C:\path\to\whatsapp-proxy-go
.\scripts\release.ps1
```

**Alternative (if execution policy blocks)**:
```powershell
powershell -ExecutionPolicy Bypass -File .\scripts\release.ps1
```

**Requirements**:
- Windows PowerShell 5.1+ or PowerShell Core 7+
- Git for Windows installed
- Git executable in PATH

---

### Cleanup Utilities

#### `cleanup-release-tag.bat` - Batch Cleanup Script
**Purpose**: Removes the `release` tag both locally and remotely.

**When to use**:
- After a successful release
- To prepare for the next release
- To fix a failed release attempt

**Usage**:
```batch
scripts\cleanup-release-tag.bat
```

---

#### `cleanup-release-tag.ps1` - PowerShell Cleanup Script
**Purpose**: PowerShell version of the cleanup utility with better output.

**Usage**:
```powershell
.\scripts\cleanup-release-tag.ps1
```

---

## 🚀 Quick Start Guide

### Creating a Release (Recommended Workflow)

#### Option 1: Using PowerShell (Recommended for Windows 10/11)

1. **Open PowerShell** in the repository root
2. **Run the release script**:
   ```powershell
   .\scripts\release.ps1
   ```
3. **Follow the prompts**:
   - Confirm you're on the correct branch
   - Review the release summary
   - Confirm the release
4. **Open GitHub Actions** to monitor the build
5. **After success, cleanup**:
   ```powershell
   .\scripts\cleanup-release-tag.ps1
   ```

#### Option 2: Using Command Prompt (Batch)

1. **Open Command Prompt** in the repository root
2. **Run the release script**:
   ```batch
   scripts\release.bat
   ```
3. **Follow the prompts** (same as PowerShell)
4. **After success, cleanup**:
   ```batch
   scripts\cleanup-release-tag.bat
   ```

#### Option 3: Manual (Advanced Users)

```bash
# Create and push tag
git tag release
git push origin release

# After workflow completes, cleanup
git tag -d release
git push origin :refs/tags/release
```

---

## 📋 Pre-Release Checklist

Before running the release script, ensure:

- [ ] All changes are committed and pushed to `main`
- [ ] All tests pass locally (`go test ./...`)
- [ ] Documentation is up to date
- [ ] CHANGELOG is updated (if applicable)
- [ ] No critical bugs in the current build
- [ ] You have push access to the repository

---

## 🔧 Troubleshooting

### "git.exe not found" Error

**Solution**:
1. Install [Git for Windows](https://git-scm.com/download/win)
2. Ensure Git is added to PATH during installation
3. Restart your terminal/command prompt

### PowerShell Execution Policy Error

**Error**: `File cannot be loaded because running scripts is disabled`

**Solution 1 - Bypass for single execution**:
```powershell
powershell -ExecutionPolicy Bypass -File .\scripts\release.ps1
```

**Solution 2 - Change policy (requires admin)**:
```powershell
Set-ExecutionPolicy RemoteSigned -Scope CurrentUser
```

### "Not a git repository" Error

**Solution**: Run the script from the repository root directory:
```batch
cd C:\path\to\whatsapp-proxy-go
scripts\release.bat
```

### Tag Already Exists

**Solution**: Run the cleanup script first:
```powershell
.\scripts\cleanup-release-tag.ps1
# Then try release again
.\scripts\release.ps1
```

### Uncommitted Changes Warning

The script will warn you if there are uncommitted changes. You can:

**Option 1 - Commit changes**:
```bash
git add .
git commit -m "your message"
git push origin main
```

**Option 2 - Stash changes**:
```bash
git stash
# Run release script
# Then restore:
git stash pop
```

**Option 3 - Continue anyway** (not recommended)
- The script will ask for confirmation

---

## 🎯 Best Practices

1. **Always use the scripts** instead of manual git commands
   - Prevents common mistakes
   - Provides safety checks
   - Shows helpful information

2. **Run from repository root** to avoid path issues

3. **Cleanup after each release** to keep tags organized

4. **Monitor GitHub Actions** to ensure successful build

5. **Test releases** by downloading and verifying binaries

---

## 📚 Additional Resources

- [Main Release Documentation](../RELEASE.md)
- [GitHub Actions Workflow](../.github/workflows/release.yml)
- [Installation Guide](../INSTALL.md)
- [Configuration Guide](../CONFIGURATION.md)

---

## 🤝 Contributing

If you find issues with these scripts or have improvements:

1. Open an issue describing the problem
2. Submit a pull request with fixes
3. Update this documentation if adding new scripts

---

## 📝 Script Maintenance

### Adding New Scripts

When adding new scripts to this directory:

1. **Use clear naming**: `verb-noun.ext` (e.g., `build-docker.bat`)
2. **Add documentation**: Update this README
3. **Include help text**: Add usage instructions in script comments
4. **Test thoroughly**: Test on clean install
5. **Handle errors**: Provide clear error messages

### Script Naming Convention

- Use lowercase with hyphens: `my-script.bat`
- Use descriptive names: `release.bat` not `r.bat`
- Include extension: `.bat`, `.ps1`, `.sh`

---

## 🔐 Security Notes

- These scripts **do not** handle credentials
- Git credentials are managed by Git credential helper
- Scripts use read-only operations except for tag creation
- Always review script contents before execution

---

## ℹ️ Version Information

- **Created**: January 2026
- **Last Updated**: January 2026
- **Maintained By**: WhatsApp Proxy Go Team

---

For questions or support, please [open an issue](https://github.com/RevEngine3r/whatsapp-proxy-go/issues) with the `scripts` label.
