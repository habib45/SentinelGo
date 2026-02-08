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

1. **Copy environment template**:
```bash
cp .env.example .env
```

2. **Edit .env file** with your Supabase credentials:
```bash
SUPABASE_URL=https://your-project.supabase.co
SUPABASE_KEY=your-supabase-anon-key
```

3. **Build the binary**:
```bash
make all
```

The Supabase credentials are now loaded from environment variables at runtime instead of being embedded in the binary.

## Configuration (Optional)

The agent requires Supabase credentials to be set via environment variables. You can optionally override defaults like heartbeat interval or GitHub repo:

### Environment Variables
- Create a `.env` file in the same directory as the binary or set system environment variables:
```bash
SUPABASE_URL=https://your-project.supabase.co
SUPABASE_KEY=your-supabase-anon-key
```

### Configuration Files
- Windows: `C:\ProgramData\sentinelgo\config.json`
- Linux: `/etc/sentinelgo/config.json`
- macOS: `/etc/sentinelgo/config.json`

Example:
```json
{
  "heartbeat_interval": "5m",
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

## Security Notes
- The agent runs as root/Administrator to collect full metrics.
- Supabase keys are loaded from environment variables at runtime; ensure the `.env` file is properly secured.
- Binary updates are fetched from GitHub Releases; ensure your repo is private or use signed releases if needed.

## License
MIT
