# SentinelGo

Cross‑platform Go agent/service that:
- Collects OS‑level metrics (CPU, memory, disk, network)
- Sends a heartbeat to Supabase (configurable interval)
- Checks for updates once per day from GitHub Releases and self‑updates
- Runs as a service on Windows, Linux (systemd), and macOS (launchd)

## Quick Start

1. **Download the binary** for your OS from [GitHub Releases](https://github.com/habib45/SentinelGo/releases/latest).
2. Follow the OS‑specific installation guide in the `doc/` folder:
   - [Windows](doc/install-windows.md)
   - [Linux](doc/install-linux.md)
   - [macOS](doc/install-macos.md)

## Build from Source

### Quick Build (Development)
```bash
# Build for current platform only
make build

# Or build all platforms
make release
```

### Release Build
```bash
# Build with specific version
make release VERSION=v1.0.0

# Build using git tag (auto-detected)
make release

# Create GitHub release (requires gh CLI)
./scripts/release.sh --release --version v1.0.0
```

### Environment Setup
1. **Copy environment template**:
```bash
cp .env.example .env
```

2. **Edit .env file** with your Supabase credentials:
```bash
SUPABASE_URL=https://your-project.supabase.co
SUPABASE_KEY=your-supabase-anon-key
API_TOKEN=your-api-token
```

The Supabase credentials are loaded from environment variables at runtime instead of being embedded in the binary.

## Configuration (Optional)

The agent requires Supabase credentials to be set via environment variables. You can optionally override defaults like heartbeat interval or GitHub repo:

### Environment Variables
- Create a `.env` file in the same directory as the binary or set system environment variables:
```bash
SUPABASE_URL=https://your-project.supabase.co
SUPABASE_KEY=your-supabase-anon-key
API_TOKEN=your-api-token
```

### Configuration Files
- Windows: `C:\ProgramData\sentinelgo\config.json`
- Linux: `/etc/sentinelgo/config.json`
- macOS: `/etc/sentinelgo/config.json`

Example:
```json
{
  "heartbeat_interval": "5m0s",
  "github_owner": "habib45",
  "github_repo": "SentinelGo",
  "current_version": "v0.1.0"
}
```

## CLI Options
```bash
./sentinelgo -install      # Install as a service (requires admin/root)
./sentinelgo -uninstall    # Uninstall the service
./sentinelgo -run          # Run in foreground (console mode)
./sentinelgo -config <path> # Use custom config file
```

## Heartbeat Payload
Sent to Supabase `/rest/v1/heartbeat`:
```json
{
  "device_id": "...",
  "version": "v0.1.0",
  "timestamp": "...",
  "alive": true,
  "system_info": { ... }
}
```

## Update Mechanism
- Every 24 hours, the agent queries GitHub Releases for the latest tag.
- If newer, it downloads the matching asset for the current OS/arch.
- It replaces the running binary and restarts.
- On Windows, a batch script handles the replace-after-exit.

## Development

### Makefile Targets
```bash
make build          # Build for current platform
make release        # Build all platforms
make test           # Run tests
make clean          # Clean build artifacts
make deps           # Download dependencies
```

### Release Process
1. **Tag the release**:
```bash
git tag v1.0.0
git push origin v1.0.0
```

2. **Automatic Release** (GitHub Actions will trigger):
   - Runs tests
   - Builds all platforms
   - Creates GitHub release with assets

3. **Manual Release** (alternative):
```bash
./scripts/release.sh --release --version v1.0.0
```

### Version Management
- Versions are injected at build time via ldflags
- Use semantic versioning (v1.0.0, v1.0.1, etc.)
- Git tags are automatically detected for versioning

## Security Notes
- The agent runs as root/Administrator to collect full metrics.
- Supabase keys are loaded from environment variables at runtime; ensure the `.env` file is properly secured.
- Binary updates are fetched from GitHub Releases; ensure your repo is private or use signed releases if needed.

## License
MIT
