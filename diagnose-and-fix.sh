#!/bin/bash

# SentinelGo Service Diagnosis and Fix Script
# Comprehensive troubleshooting for service issues

set -e

echo "🔍 SentinelGo Service Diagnosis & Fix"
echo "=================================="

# 1. Check current service status
echo "📊 Current Service Status:"
sudo systemctl status sentinelgo --no-pager -l

echo ""
echo "📋 Recent Service Logs (last 10 entries):"
sudo journalctl -u sentinelgo -n 10 --no-pager

echo ""
echo "🔍 Checking Common Issues..."

# 2. Check if binary exists and is executable
echo "📁 Binary Status:"
if [[ -x "/opt/sentinelgo/sentinelgo" ]]; then
    echo "✅ Binary exists and is executable"
else
    echo "❌ Binary missing or not executable"
    echo "🔧 Fixing permissions..."
    sudo chmod +x /opt/sentinelgo/sentinelgo
fi

# 3. Check config file
echo ""
echo "📝 Configuration Status:"
if [[ -f "/opt/sentinelgo/.sentinelgo/config.json" ]]; then
    echo "✅ Config file exists"
    echo "📄 Config contents:"
    cat "/opt/sentinelgo/.sentinelgo/config.json"
else
    echo "❌ Config file missing"
    echo "🔧 Creating default config..."
    sudo mkdir -p "/opt/sentinelgo/.sentinelgo"
    echo '{"heartbeat_interval":"5m0s","auto_update":false}' | sudo tee "/opt/sentinelgo/.sentinelgo/config.json"
    echo "✅ Config file created"
fi

# 4. Check permissions
echo ""
echo "🔐 Permission Status:"
echo "Directory: /opt/sentinelgo"
ls -la /opt/sentinelgo | grep -E "(sentinelgo|\.sentinelgo)"

echo ""
echo "User: $SERVICE_USER"
echo "Group: $SERVICE_USER"

# 5. Test binary manually
echo ""
echo "🧪 Manual Binary Test:"
echo "Testing if binary can run at all..."
sudo -u sentinelgo /opt/sentinelgo/sentinelgo --version 2>&1

# 6. Check dependencies
echo ""
echo "📦 Dependency Check:"
echo "Checking for missing libraries..."
ldd /opt/sentinelgo/sentinelgo 2>/dev/null | grep "not found" || echo "✅ All dependencies found"

# 7. Environment check
echo ""
echo "🌍 Environment Variables:"
echo "HOME: $(sudo -u sentinelgo printenv HOME 2>/dev/null || echo 'Not set')"
echo "USER: $(sudo -u sentinelgo printenv USER 2>/dev/null || echo 'Not set')"

# 8. Port check (if applicable)
echo ""
echo "🔌 Port Check:"
echo "Checking if port 8080 is available..."
if command -v netstat >/dev/null 2>&1; then
    if sudo -u sentinelgo netstat -tlnp 2>/dev/null | grep -q ":8080"; then
        echo "⚠️  Port 8080 appears to be in use"
    else
        echo "✅ Port 8080 is available"
    fi
else
    echo "ℹ️  netstat not available - skipping port check"
fi

# 9. Process check
echo ""
echo "🔍 Process Check:"
echo "Checking for running SentinelGo processes..."
ps aux | grep -i sentinelgo | grep -v grep || echo "✅ No unexpected processes found"

# 10. Memory and disk space
echo ""
echo "💾 System Resources:"
echo "Memory usage:"
free -h | grep -E "^Mem:" | awk '{print "  Total: " $2 " Used: " $3 " Free: " $4}'
echo ""
echo "Disk space:"
df -h /opt/sentinelgo | tail -1

echo ""
echo "🔧 Recommended Fixes:"

# Based on diagnosis, provide specific fixes
if [[ ! -x "/opt/sentinelgo/sentinelgo" ]]; then
    echo "1. Fix binary permissions"
    echo "   sudo chmod +x /opt/sentinelgo/sentinelgo"
fi

if [[ ! -f "/opt/sentinelgo/.sentinelgo/config.json" ]]; then
    echo "2. Fix missing config"
    echo "   The config creation above should fix this"
fi

# Check if service is in restart loop
RESTART_COUNT=$(journalctl -u sentinelgo | grep -c "Started sentinelgo.service" | wc -l)
if [[ $RESTART_COUNT -gt 3 ]]; then
    echo "3. Service restart loop detected"
    echo "   sudo systemctl stop sentinelgo"
    echo "   sudo systemctl reset-failed sentinelgo"
    echo "   sleep 5"
    echo "   sudo systemctl start sentinelgo"
fi

echo ""
echo "🎯 Next Steps:"
echo "1. If issues persist, try: sudo ./install.sh fix-service"
echo "2. Check logs: sudo journalctl -u sentinelgo -f"
echo "3. Test manually: sudo -u sentinelgo /opt/sentinelgo/sentinelgo -run"
echo "4. Enable auto-updates: sudo ./sentinelgo -enable-auto-update"

echo ""
echo "✅ Diagnosis complete!"
