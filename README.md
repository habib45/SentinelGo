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

## Build from Source (with Supabase credentials)

Set environment variables before building:
```bash
export SUPABASE_URL=https://myproject.supabase.co
export SUPABASE_KEY=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
make all
```
Or manually:
```bash
go build -ldflags "-X sentinelgo/internal/heartbeat.SupabaseURL=$SUPABASE_URL -X sentinelgo/internal/heartbeat.SupabaseKey=$SUPABASE_KEY" ./cmd/sentinelgo
```

## Configuration (Optional)

The agent works out of the box with build‑time embedded Supabase credentials.
You can optionally override defaults like heartbeat interval or GitHub repo:

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
- Supabase keys are embedded at build time; users cannot change them.
- Binary updates are fetched from GitHub Releases; ensure your repo is private or use signed releases if needed.

## License
MIT
