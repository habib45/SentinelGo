package osinfo

import (
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
	Timestamp   time.Time  `json:"timestamp"`
	Hostname    string     `json:"hostname"`
	OS          string     `json:"os"`
	Platform    string     `json:"platform"`
	PlatformVer string     `json:"platform_version"`
	Arch        string     `json:"arch"`
	Uptime      uint64     `json:"uptime"`
	CPU         CPUInfo    `json:"cpu"`
	Memory      MemoryInfo `json:"memory"`
	Disk        DiskInfo   `json:"disk"`
	Network     []NetInfo  `json:"network"`
	EmployeeId  string     `json:"employee_id"`
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
}

func Collect() *SystemInfo {
	hInfo, _ := host.Info()
	cpuInfo, _ := cpu.Info()
	cpuPercent, _ := cpu.Percent(0, false)
	memInfo, _ := mem.VirtualMemory()
	diskInfo, _ := disk.Usage("/")
	netInfo, _ := net.IOCounters(true)

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
		netStats = append(netStats, NetInfo{
			Name:      ni.Name,
			BytesSent: ni.BytesSent,
			BytesRecv: ni.BytesRecv,
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
	if runtime.GOOS == "darwin" {
		osName = "IOS"
	} else if runtime.GOOS == "linux" {
		osName = "linux"
	} else if runtime.GOOS == "windows" {
		osName = "windows"
	}

	return &SystemInfo{
		Timestamp:   time.Now(),
		Hostname:    hInfo.Hostname,
		OS:          osName,
		Platform:    hInfo.Platform,
		PlatformVer: hInfo.PlatformVersion,
		Arch:        runtime.GOARCH,
		Uptime:      hInfo.Uptime,
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
	}
}
