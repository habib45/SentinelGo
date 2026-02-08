# Release Guide

This document explains how to create releases for SentinelGo.

## Quick Release Process

### 1. Automated Release (Recommended)
```bash
# Create and push a git tag
git tag v1.0.0
git push origin v1.0.0

# GitHub Actions will automatically:
# - Run tests
# - Build all platforms
# - Create GitHub release with assets
```

### 2. Manual Release
```bash
# Build and release manually
./scripts/release.sh --release --version v1.0.0
```

## Build Commands

### Development Build
```bash
# Build for current platform only
make build

# Build all platforms (no version specified)
make release

# Build with specific version
make release VERSION=v1.0.0
```

### Release Assets
The build process creates these assets:
- `sentinelgo-windows-amd64.exe` (~8MB)
- `sentinelgo-linux-amd64` (~9.8MB)
- `sentinelgo-linux-arm64` (~9.5MB)
- `sentinelgo-darwin-amd64` (~9.7MB)
- `sentinelgo-darwin-arm64` (~9.4MB)

## Version Management

### Version Sources (in order of priority)
1. `VERSION` environment variable
2. Git tag (auto-detected)
3. Falls back to "dev"

### Version Injection
Versions are injected at build time via ldflags:
```bash
go build -ldflags "-X sentinelgo/internal/config.Version=$VERSION"
```

### Semantic Versioning
Use semantic versioning for releases:
- `v1.0.0` - Major release
- `v1.0.1` - Patch release
- `v1.1.0` - Minor release

## Environment Variables

Users must configure these environment variables:
```bash
SUPABASE_URL=https://your-project.supabase.co
SUPABASE_KEY=your-supabase-anon-key
API_TOKEN=your-api-token
```

## GitHub Actions

### CI Workflow (.github/workflows/ci.yml)
- Runs on push to main/develop branches
- Runs tests on pull requests
- Builds all platforms for testing

### Release Workflow (.github/workflows/release.yml)
- Triggers on git tags (v*)
- Runs full test suite
- Builds all platforms
- Creates GitHub release with assets

## Security Considerations

- No secrets embedded in binaries
- Environment variables loaded at runtime
- GitHub releases are public by default
- Ensure Supabase credentials are properly secured

## Troubleshooting

### Build Issues
```bash
# Clean build artifacts
make clean

# Rebuild dependencies
make deps

# Check Go version (requires 1.22+)
go version
```

### Release Issues
```bash
# Check GitHub CLI installation
gh --version

# Verify authentication
gh auth status

# Manual release creation
gh release create v1.0.0 build/*
```

### Environment Issues
```bash
# Check .env file exists
ls -la .env

# Verify environment variables
cat .env
```

## Migration from Old Build System

The old build system embedded Supabase credentials at build time. The new system:

1. **Removes embedded secrets** - Better security
2. **Uses environment variables** - Runtime configuration
3. **Injects version only** - Simpler build process
4. **Automates releases** - GitHub Actions integration

No changes needed for existing installations - they'll continue to work with the new binary format.
