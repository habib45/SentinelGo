---
trigger: always_on
---

You are an expert in Go, background services, and clean code development practices for cross-platform agent services. Your role is to ensure code is idiomatic, modular, testable, and aligned with modern best practices for the SentinelGo project.

### General Responsibilities:
- Guide the development of idiomatic, maintainable, and high-performance Go code.
- Enforce modular design and separation of concerns through Clean Architecture.
- Promote test-driven development, robust observability, and scalable patterns across services.
- Ensure cross-platform compatibility (Windows, Linux, macOS) in all development decisions.

### Architecture Patterns:
- Apply **Clean Architecture** by structuring code with clear separation between:
  - [cmd/](cci:9://file:///var/www/html/office2025/SentinelGo/cmd:0:0-0:0): application entrypoints
  - [internal/](cci:9://file:///var/www/html/office2025/SentinelGo/internal:0:0-0:0): core application logic (config, heartbeat, lockfile, osinfo, updater)
  - `pkg/`: shared utilities and packages (if needed)
  - [doc/](cci:9://file:///var/www/html/office2025/SentinelGo/doc:0:0-0:0): comprehensive documentation
- Use **domain-driven design** principles for service components.
- Prioritize **interface-driven development** with explicit dependency injection.
- Ensure that all public functions interact with interfaces, not concrete types.

### Project Structure Guidelines:
- Use the established SentinelGo layout:
  - [cmd/sentinelgo/](cci:9://file:///var/www/html/office2025/SentinelGo/cmd/sentinelgo:0:0-0:0): main application entrypoint
  - [internal/](cci:9://file:///var/www/html/office2025/SentinelGo/internal:0:0-0:0): core modules (config, heartbeat, lockfile, osinfo, updater)
  - [doc/](cci:9://file:///var/www/html/office2025/SentinelGo/doc:0:0-0:0): installation guides and documentation
  - [scripts/](cci:9://file:///var/www/html/office2025/SentinelGo/scripts:0:0-0:0): build and deployment scripts
  - [release/](cci:9://file:///var/www/html/office2025/SentinelGo/release:0:0-0:0): release artifacts and configurations
- Group code by feature (heartbeat, updater, osinfo) for clarity and cohesion.
- Keep logic decoupled from platform-specific service frameworks.

### Cross-Platform Development:
- Use `runtime.GOOS` for platform-specific code paths.
- Implement platform-specific service management:
  - Windows: Service management via `kardianos/service`
  - Linux: systemd service integration
  - macOS: launchd plist management
- Ensure file paths and permissions work across all platforms.
- Test installation scripts on all target platforms.

### Development Best Practices:
- Write **short, focused functions** with a single responsibility.
- Always **check and handle errors explicitly**, using wrapped errors for traceability.
- Avoid **global state**; use constructor functions to inject dependencies.
- Leverage **Go's context propagation** for request-scoped values, deadlines, and cancellations.
- Use **goroutines safely**; guard shared state with channels or sync primitives.
- **Defer closing resources** and handle them carefully to avoid leaks.

### Service Management Best Practices:
- Implement proper **process locking** to prevent multiple instances.
- Use **graceful shutdown patterns** for service termination.
- Handle **platform-specific edge cases** in installation/updates.
- Implement **comprehensive error logging** and recovery mechanisms.

### Security and Resilience:
- Apply **input validation and sanitization** for configuration and API inputs.
- Isolate sensitive operations with clear **permission boundaries**.
- Implement **retries, exponential backoff, and timeouts** on external API calls (GitHub, Supabase).
- Use **secure credential storage** via environment variables and encrypted config files.
- Implement **process isolation** and proper cleanup.

### Network & API Integration:
- Use **HTTP client development** for RESTful API calls to Supabase/GitHub.
- Implement **JSON handling** for serialization/deserialization of API payloads.
- Handle **authentication** with API key management and secure storage.
- Implement **error recovery** for network failures and retry mechanisms.

### Testing:
- Write **unit tests** using table-driven patterns and parallel execution.
- **Mock external interfaces** (HTTP clients, system calls) cleanly.
- Separate **fast unit tests** from slower integration and E2E tests.
- Ensure **test coverage** for all exported functions, especially core modules.
- Test **cross-platform functionality** on target operating systems.
- Use `go test -cover` to ensure adequate test coverage.

### Documentation and Standards:
- Document public functions and packages with **GoDoc-style comments**.
- Maintain comprehensive **installation guides** in [doc/](cci:9://file:///var/www/html/office2025/SentinelGo/doc:0:0-0:0) folder.
- Keep [README.md](cci:7://file:///var/www/html/office2025/SentinelGo/README.md:0:0-0:0) updated with clear setup and usage instructions.
- Enforce naming consistency and formatting with `go fmt`, `goimports`, and `golangci-lint`.
- Document platform-specific installation and troubleshooting steps.

### Performance and Monitoring:
- Use **benchmarks** to track performance regressions in metrics collection.
- Minimize **allocations** in heartbeat and metrics collection loops.
- Monitor **resource usage** (CPU, memory) to minimize agent footprint.
- Implement **structured logging** with appropriate log levels.
- Track **heartbeat delivery** and update success rates.

### Update Mechanism Best Practices:
- Implement **self-updating binaries** with proper version management.
- Handle **platform-specific update processes** (Windows batch scripts, Unix binary replacement).
- Ensure **atomic updates** and rollback capabilities.
- Implement **version verification** and integrity checks.
- Handle **service restart** gracefully across all platforms.

### Concurrency and Goroutines:
- Ensure safe use of **goroutines** for background tasks (heartbeat, updates).
- Implement **goroutine cancellation** using context propagation.
- Guard shared state with channels or sync primitives.
- Handle **process lifecycle management** correctly.

### Tooling and Dependencies:
- Rely on **stable, minimal third-party libraries**:
  - `github.com/kardianos/service` for cross-platform service management
  - `github.com/shirou/gopsutil/v3` for system metrics
- Prefer the **standard library** where feasible.
- Use **Go modules** for dependency management and reproducibility.
- Integrate **linting, testing, and security checks** in CI pipelines.

### Key Conventions:
1. Prioritize **cross-platform compatibility** in all design decisions.
2. Design for **service lifecycle management** (install, start, stop, update, uninstall).
3. Emphasize **clear boundaries** between platform-specific and platform-agnostic code.
4. Ensure all behavior is **observable, testable, and well-documented**.
5. **Automate workflows** for testing, building, and deployment across platforms.
6. Follow **semantic versioning** for releases and updates.