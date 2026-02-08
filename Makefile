# SentinelGo Makefile
# Build with Supabase credentials embedded at build time

# Set these before building, e.g.:
# export SUPABASE_URL=https://myproject.supabase.co
# export SUPABASE_KEY=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

# Build flags
LDFLAGS=-ldflags "-X sentinelgo/internal/heartbeat.SupabaseURL=$(SUPABASE_URL) -X sentinelgo/internal/heartbeat.SupabaseKey=$(SUPABASE_KEY)"

# Targets
.PHONY: build clean all windows linux macos

all: windows linux macos

windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o build/windows/sentinelgo-windows-amd64.exe ./cmd/sentinelgo

linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o build/linux/sentinelgo-linux-amd64 ./cmd/sentinelgo
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o build/linux/sentinelgo-linux-arm64 ./cmd/sentinelgo

macos:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o build/darwin/sentinelgo-darwin-amd64 ./cmd/sentinelgo
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o build/darwin/sentinelgo-darwin-arm64 ./cmd/sentinelgo

clean:
	rm -rf build/

# Example usage:
# export SUPABASE_URL=https://myproject.supabase.co
# export SUPABASE_KEY=your-anon-key
# make all
