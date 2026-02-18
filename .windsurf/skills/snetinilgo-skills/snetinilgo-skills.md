---
name: snetinilgo-skills
description: Comprehensive skills for developing cross-platform Go agent services that run on Windows, Linux, and macOS with auto-start functionality
---
 
# Cross-Platform Go Agent Development Skills
 
## Core Go Development Skills
- **Go Programming Fundamentals** - Strong understanding of Go syntax, packages, and project structure
- **Cross-Platform Development** - Experience with `runtime.GOOS` for OS-specific code paths
- **Concurrency** - Goroutines and channels for background tasks (heartbeat, updates)
- **Error Handling** - Proper Go error handling patterns
- **File System Operations** - Cross-platform file operations and permissions
 
## Service Management Skills
- **Windows Service Development** - Using `github.com/kardianos/service` for Windows services
- **Linux systemd Services** - Understanding systemd service files and management
- **macOS launchd Services** - Creating and managing launchd plist files
- **Process Management** - Cross-platform process detection and termination
 
## System Integration Skills
- **OS-Specific Metrics Collection** - CPU, memory, disk, network monitoring using `gopsutil`
- **System API Integration** - Windows APIs, Linux procfs, macOS system calls
- **Auto-Start Configuration** - Registry entries (Windows), systemd services (Linux), launchd (macOS)
- **Permission Management** - Handling administrator/root privileges
 
## Network & API Skills
- **HTTP Client Development** - RESTful API calls to Supabase/GitHub
- **JSON Handling** - Serialization/deserialization of API payloads
- **Authentication** - API key management and secure credential storage
- **Error Recovery** - Network failure handling and retry mechanisms
 
## DevOps & Deployment Skills
- **Cross-Platform Build Systems** - Go build tags and conditional compilation
- **Binary Distribution** - Creating platform-specific executables
- **Installation Scripts** - Batch files (Windows), shell scripts (Linux/macOS)
- **Update Mechanisms** - Self-updating binaries and version management
 
## Security Skills
- **Secure Credential Storage** - Environment variables and encrypted config files
- **Process Isolation** - Preventing multiple instances and proper cleanup
- **Code Signing** - Binary verification and secure updates
- **Access Control** - Proper file permissions and user privileges
 
## Monitoring & Logging Skills
- **Structured Logging** - Cross-platform logging with different levels
- **System Metrics Collection** - Performance monitoring and health checks
- **Heartbeat Implementation** - Regular status reporting to backend
- **Debugging Tools** - Cross-platform debugging and troubleshooting
 
## Configuration Management Skills
- **Environment Variable Handling** - Cross-platform config loading
- **JSON Configuration** - Persistent settings and user preferences
- **Default Value Management** - Sensible defaults and migration handling
- **Validation** - Input validation and error reporting
 
## Key Libraries & Tools
- `github.com/kardianos/service` - Cross-platform service management
- `github.com/shirou/gopsutil/v3` - System metrics collection
- Standard library packages: `os`, `os/exec`, `runtime`, `net/http`, `encoding/json`
 
## Best Practices
- Implement proper process locking to prevent multiple instances
- Use graceful shutdown patterns for service termination
- Handle platform-specific edge cases in installation/updates
- Implement comprehensive error logging and recovery
- Follow semantic versioning for releases
- Test thoroughly on all target platforms