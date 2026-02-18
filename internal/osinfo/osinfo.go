package osinfo

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

type SystemInfo struct {
	Timestamp       time.Time  `json:"timestamp"`
	Hostname        string     `json:"hostname"`
	OS              string     `json:"os"`
	Platform        string     `json:"platform"`
	PlatformVer     string     `json:"platform_version"`
	Arch            string     `json:"arch"`
	Uptime          uint64     `json:"uptime"`
	UptimeFormatted string     `json:"uptime_formatted"`
	CPU             CPUInfo    `json:"cpu"`
	Memory          MemoryInfo `json:"memory"`
	Disk            DiskInfo   `json:"disk"`
	Network         []NetInfo  `json:"network"`
	EmployeeId      string     `json:"employee_id"`
	MACAddress      string     `json:"mac_address"`
}

type CPUInfo struct {
	ModelName string  `json:"model_name"`
	Cores     int     `json:"cores"`
	Usage     float64 `json:"usage_percent"`
}

type MemoryInfo struct {
	Total uint64  `json:"total"`
	Used  uint64  `json:"used"`
	Free  uint64  `json:"free"`
	Usage float64 `json:"usage_percent"`
}

type DiskInfo struct {
	Total uint64 `json:"total"`
	Used  uint64 `json:"used"`
	Free  uint64 `json:"free"`
}

type NetInfo struct {
	Name      string `json:"name"`
	BytesSent uint64 `json:"bytes_sent"`
	BytesRecv uint64 `json:"bytes_recv"`
	MACAddr   string `json:"mac_address"`
}

// formatUptime converts uptime in seconds to human-readable format
func formatUptime(seconds uint64) string {
	if seconds == 0 {
		return "0 seconds"
	}

	days := seconds / 86400
	hours := (seconds % 86400) / 3600
	minutes := (seconds % 3600) / 60
	secs := seconds % 60

	var parts []string

	if days > 0 {
		if days == 1 {
			parts = append(parts, "1 day")
		} else {
			parts = append(parts, fmt.Sprintf("%d days", days))
		}
	}

	if hours > 0 {
		if hours == 1 {
			parts = append(parts, "1 hour")
		} else {
			parts = append(parts, fmt.Sprintf("%d hours", hours))
		}
	}

	if minutes > 0 {
		if minutes == 1 {
			parts = append(parts, "1 minute")
		} else {
			parts = append(parts, fmt.Sprintf("%d minutes", minutes))
		}
	}

	// Only show seconds if no other time units or if uptime is very short
	if len(parts) == 0 || secs > 0 {
		if secs == 1 {
			parts = append(parts, "1 second")
		} else {
			parts = append(parts, fmt.Sprintf("%d seconds", secs))
		}
	}

	// Join with spaces, limit to 3 most significant units
	if len(parts) > 3 {
		parts = parts[:3]
	}

	return strings.Join(parts, " ")
}

// getPrimaryMACAddress returns the first non-empty MAC address from network interfaces
func getPrimaryMACAddress(netInterfaces []net.InterfaceStat) string {
	for _, iface := range netInterfaces {
		if iface.HardwareAddr != "" && !strings.HasPrefix(iface.HardwareAddr, "00:00:00") {
			return iface.HardwareAddr
		}
	}
	return ""
}

func Collect() *SystemInfo {
	hInfo, _ := host.Info()
	cpuInfo, _ := cpu.Info()
	cpuPercent, _ := cpu.Percent(0, false)
	memInfo, _ := mem.VirtualMemory()
	diskInfo, _ := disk.Usage("/")
	netInfo, _ := net.IOCounters(true)
	netInterfaces, _ := net.Interfaces()

	var cpuModel string
	if len(cpuInfo) > 0 {
		cpuModel = cpuInfo[0].ModelName
	}

	var cpuUsage float64
	if len(cpuPercent) > 0 {
		cpuUsage = cpuPercent[0]
	}

	var netStats []NetInfo
	for _, ni := range netInfo {
		// Find MAC address for this interface
		var macAddr string
		for _, iface := range netInterfaces {
			if iface.Name == ni.Name {
				macAddr = iface.HardwareAddr
				break
			}
		}

		netStats = append(netStats, NetInfo{
			Name:      ni.Name,
			BytesSent: ni.BytesSent,
			BytesRecv: ni.BytesRecv,
			MACAddr:   macAddr,
		})
	}

	// Clean up hostname for macOS (remove .local suffix)
	hostname := hInfo.Hostname
	if runtime.GOOS == "darwin" {
		hostname = strings.TrimSuffix(hostname, ".local")
		hostname = strings.TrimSuffix(hostname, ".lan")
		hostname = strings.TrimSuffix(hostname, ".home")
	}

	// Normalize OS name for better readability
	osName := hInfo.OS
	switch runtime.GOOS {
	case "darwin":
		osName = "macOS"
	case "linux":
		osName = "linux"
	case "windows":
		osName = "windows"
	}

	return &SystemInfo{
		Timestamp:       time.Now(),
		Hostname:        hInfo.Hostname,
		OS:              osName,
		Platform:        hInfo.Platform,
		PlatformVer:     hInfo.PlatformVersion,
		Arch:            runtime.GOARCH,
		Uptime:          hInfo.Uptime,
		UptimeFormatted: formatUptime(hInfo.Uptime),
		CPU: CPUInfo{
			ModelName: cpuModel,
			Cores:     runtime.NumCPU(),
			Usage:     cpuUsage,
		},
		Memory: MemoryInfo{
			Total: memInfo.Total,
			Used:  memInfo.Used,
			Free:  memInfo.Free,
			Usage: memInfo.UsedPercent,
		},
		Disk: DiskInfo{
			Total: diskInfo.Total,
			Used:  diskInfo.Used,
			Free:  diskInfo.Free,
		},
		Network:    netStats,
		EmployeeId: hostname,
		MACAddress: getPrimaryMACAddress(netInterfaces),
	}
}
