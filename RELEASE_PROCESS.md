# SentinelGo Release Process

## 🚀 **Automated Release Workflow**

This document explains how to create automated releases for SentinelGo using GitHub Actions.

---

## 📋 **Prerequisites**

1. **Git repository setup** with GitHub Actions enabled
2. **GitHub token** with repository write permissions (automatically provided by Actions)
3. **Clean working directory** - no uncommitted changes
4. **Main branch** - releases must be created from the main branch

---

## 🎯 **Quick Release Process**

### **Option 1: Using the Release Script (Recommended)**

```bash
# Create a new release
./release.sh v1.0.0

# The script will:
# 1. Build and test locally
# 2. Show what packages will be created
# 3. Ask for confirmation
# 4. Create and push the tag
# 5. Trigger GitHub Actions
```

### **Option 2: Manual Tag Creation**

```bash
# Create and push tag manually
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

---

## 🔄 **What Happens Automatically**

When you push a tag (starting with `v`), GitHub Actions will:

### **1. Build Phase**
- ✅ **Test** - Run all unit tests
- ✅ **Build** - Create binaries for all platforms:
  - `sentinelgo-linux-amd64`
  - `sentinelgo-linux-arm64`
  - `sentinelgo-darwin-amd64`
  - `sentinelgo-darwin-arm64`
  - `sentinelgo-windows-amd64.exe`

### **2. Package Creation**
- ✅ **Package** - Create tar.gz files with:
  - Platform-specific binary
  - Universal installer (`install.sh`)
  - Documentation (`INSTALLATION.md`)

### **3. Release Creation**
- ✅ **Draft Release** - Create draft release on GitHub
- ✅ **Release Notes** - Generate automatically from git commits
- ✅ **Upload Assets** - Attach all package files

---

## 📦 **Generated Release Packages**

```
sentinelgo-v1.0.0-linux-amd64.tar.gz     (5.8MB)
sentinelgo-v1.0.0-linux-arm64.tar.gz     (5.4MB)
sentinelgo-v1.0.0-darwin-amd64.tar.gz     (4.7MB)
sentinelgo-v1.0.0-darwin-arm64.tar.gz     (4.4MB)
sentinelgo-v1.0.0-windows.tar.gz          (4.8MB)
```

Each package contains:
- **Binary** for the specific platform/architecture
- **Universal installer** (`install.sh`)
- **Complete documentation** (`INSTALLATION.md`)

---

## 🎛️ **Manual Release Publishing**

After GitHub Actions creates the draft release:

1. **Go to GitHub Releases**
   - Navigate to your repository
   - Click "Releases" tab
   - Find the draft release

2. **Review the Release**
   - Check the automatically generated release notes
   - Verify all package files are attached
   - Make any necessary edits

3. **Publish the Release**
   - Click "Publish release"
   - The release becomes public and downloadable

---

## 📋 **Release Notes Template**

The workflow automatically generates release notes including:

### **Changelog**
- Commits since the last tag
- Formatted as bullet points

### **Installation Instructions**
- Quick install guide
- Available packages list
- Feature highlights
- Documentation links

### **Example Output**
```markdown
## Changelog
- Add automatic update functionality
- Fix duplicate heartbeat issue
- Improve cross-platform compatibility

## 🚀 Installation

### Quick Install (Recommended)
1. Download the package for your OS/architecture below
2. Extract the archive
3. Run: sudo ./install.sh
4. Enable auto-updates: sudo ./sentinelgo -enable-auto-update

### Available Packages
- Linux AMD64: For standard Ubuntu/Debian/CentOS systems
- Linux ARM64: For ARM systems like Raspberry Pi
- macOS AMD64: For Intel Macs
- macOS ARM64: For Apple Silicon Macs (M1/M2)
- Windows: For Windows 10/11 systems
```

---

## 🔧 **Configuration**

### **GitHub Workflow Settings**

The workflow is configured in `.github/workflows/release.yml`:

- **Trigger**: Tags starting with `v` (e.g., `v1.0.0`)
- **Go Version**: 1.22
- **Platforms**: Linux, macOS, Windows (AMD64 + ARM64)
- **Release Type**: Draft (requires manual publishing)

### **Environment Variables**

The workflow uses these environment variables:
- `GO_VERSION`: Go compiler version
- `GITHUB_TOKEN`: Automatically provided by GitHub Actions

---

## 🛠️ **Troubleshooting**

### **Common Issues**

#### **Workflow Fails**
```bash
# Check workflow logs on GitHub
# Common causes:
# - Build errors in code
# - Test failures
# - GitHub Actions service issues
```

#### **Release Not Created**
```bash
# Verify tag was pushed
git tag -l
git ls-remote --tags origin

# Check if tag follows naming convention
# Must start with 'v' (e.g., v1.0.0)
```

#### **Packages Missing**
```bash
# Check workflow artifacts
# Verify build completed successfully
# Check file paths in workflow
```

### **Manual Recovery**

If the automated process fails:

1. **Build locally**
   ```bash
   make clean
   make release
   ```

2. **Create manual release**
   - Go to GitHub Releases
   - Click "Create new release"
   - Upload packages from `release/` directory
   - Write release notes manually

---

## 🎯 **Best Practices**

### **Version Naming**
- Use semantic versioning: `v1.0.0`, `v1.0.1`, `v1.1.0`
- Always start with `v`
- Use consistent format

### **Pre-Release Checklist**
- [ ] All tests pass locally
- [ ] Documentation is updated
- [ ] Version is bumped in code
- [ ] Working directory is clean
- [ ] On main branch

### **Post-Release**
- [ ] Verify release is published
- [ ] Test installation from packages
- [ ] Update changelog if needed
- [ ] Announce release to users

---

## 📞 **Support**

### **Getting Help**
- Check GitHub Actions logs for workflow issues
- Review this documentation for common problems
- Test locally before creating releases

### **Emergency Rollback**
If a release has issues:
1. **Delete the release** on GitHub
2. **Create a new tag** with fixed version
3. **Publish new release**

---

## 🎉 **Success!**

Once completed, you'll have:
- ✅ **Automated builds** for all platforms
- ✅ **Professional packages** with installers
- ✅ **Draft releases** for review
- ✅ **Complete documentation** for users
- ✅ **Easy deployment** for employees

Your SentinelGo is ready for enterprise deployment! 🚀
