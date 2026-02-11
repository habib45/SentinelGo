package lockfile

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// LockFile represents a process lock file
type LockFile struct {
	path     string
	file     *os.File
	acquired bool
}

// NewLockFile creates a new lock file instance
func NewLockFile(name string) *LockFile {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "/tmp"
	}
	lockDir := filepath.Join(home, ".sentinelgo")
	os.MkdirAll(lockDir, 0755)
	
	return &LockFile{
		path: filepath.Join(lockDir, name+".lock"),
	}
}

// TryAcquire attempts to acquire the lock
func (lf *LockFile) TryAcquire() error {
	// Try to create/open the lock file
	file, err := os.OpenFile(lf.path, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("open lock file: %w", err)
	}
	lf.file = file

	// Try to acquire exclusive lock using syscall
	err = syscall.Flock(int(file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		file.Close()
		lf.file = nil
		if err == syscall.EWOULDBLOCK {
			// Lock is held by another process
			return fmt.Errorf("lock already held by another process")
		}
		return fmt.Errorf("flock failed: %w", err)
	}

	// Write our PID to the lock file
	pid := os.Getpid()
	_, err = file.WriteString(strconv.Itoa(pid) + "\n")
	if err != nil {
		syscall.Flock(int(file.Fd()), syscall.LOCK_UN)
		file.Close()
		lf.file = nil
		return fmt.Errorf("write PID: %w", err)
	}

	// Sync to ensure PID is written to disk
	file.Sync()

	lf.acquired = true
	return nil
}

// AcquireWithTimeout attempts to acquire the lock with a timeout
func (lf *LockFile) AcquireWithTimeout(timeout time.Duration) error {
	start := time.Now()
	for time.Since(start) < timeout {
		err := lf.TryAcquire()
		if err == nil {
			return nil
		}
		if !strings.Contains(err.Error(), "lock already held") {
			return err
		}
		time.Sleep(100 * time.Millisecond)
	}
	return fmt.Errorf("timeout acquiring lock after %v", timeout)
}

// Release releases the lock
func (lf *LockFile) Release() error {
	if !lf.acquired || lf.file == nil {
		return nil
	}

	// Remove the lock file
	os.Remove(lf.path)

	// Release the flock
	err := syscall.Flock(int(lf.file.Fd()), syscall.LOCK_UN)
	if err != nil {
		return fmt.Errorf("unlock failed: %w", err)
	}

	// Close the file
	lf.file.Close()
	lf.file = nil
	lf.acquired = false

	return nil
}

// GetLockedPID returns the PID of the process holding the lock
func (lf *LockFile) GetLockedPID() (int, error) {
	data, err := os.ReadFile(lf.path)
	if err != nil {
		return 0, err
	}

	pidStr := strings.TrimSpace(string(data))
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		return 0, fmt.Errorf("invalid PID in lock file: %s", pidStr)
	}

	return pid, nil
}

// IsProcessRunning checks if a process with the given PID is running
func IsProcessRunning(pid int) bool {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	// Send signal 0 to check if process exists
	err = process.Signal(syscall.Signal(0))
	return err == nil
}

// CheckExistingLock checks if there's an existing lock and if the process is still running
func (lf *LockFile) CheckExistingLock() (bool, error) {
	pid, err := lf.GetLockedPID()
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil // No lock file exists
		}
		return false, err
	}

	// Check if the process is still running
	if IsProcessRunning(pid) {
		return true, nil // Process is still running
	}

	// Process is dead, clean up the stale lock file
	os.Remove(lf.path)
	return false, nil
}
