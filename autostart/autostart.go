// Package autostart provides functionality to configure the application to start automatically on system boot
package autostart

import (
	"errors"
	"fmt"
	"os"
	"runtime"
)

// Common errors
var (
	ErrUnsupportedOS = errors.New("unsupported operating system")
)

// AutoStart handles system autostart configuration
type AutoStart struct {
	Name        string // Application name
	DisplayName string // Display name for system entries
	Executable  string // Path to executable
}

// New creates a new AutoStart instance for the specified application
func New(name, displayName string) *AutoStart {
	executable, _ := os.Executable()
	return &AutoStart{
		Name:        name,
		DisplayName: displayName,
		Executable:  executable,
	}
}

// IsEnabled checks if the application is configured to start automatically
func (a *AutoStart) IsEnabled() bool {
	switch runtime.GOOS {
	case "windows":
		return a.isEnabledWindows()
	case "linux":
		return a.isEnabledLinux()
	case "darwin":
		return a.isEnabledMacOS()
	default:
		return false
	}
}

// Enable configures the application to start automatically on system boot
func (a *AutoStart) Enable() error {
	switch runtime.GOOS {
	case "windows":
		return a.enableWindows()
	case "linux":
		return a.enableLinux()
	case "darwin":
		return a.enableMacOS()
	default:
		return ErrUnsupportedOS
	}
}

// Disable removes the configuration for the application to start automatically
func (a *AutoStart) Disable() error {
	switch runtime.GOOS {
	case "windows":
		return a.disableWindows()
	case "linux":
		return a.disableLinux()
	case "darwin":
		return a.disableMacOS()
	default:
		return ErrUnsupportedOS
	}
}

// Windows implementation
func (a *AutoStart) isEnabledWindows() bool {
	// TODO: Check registry for startup entry
	return false
}

func (a *AutoStart) enableWindows() error {
	// TODO: Create registry entry in HKCU\Software\Microsoft\Windows\CurrentVersion\Run
	return nil
}

func (a *AutoStart) disableWindows() error {
	// TODO: Remove registry entry
	return nil
}

// Linux implementation
func (a *AutoStart) isEnabledLinux() bool {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false
	}
	path := fmt.Sprintf("%s/.config/autostart/%s.desktop", homeDir, a.Name)
	_, err = os.Stat(path)
	return err == nil
}

func (a *AutoStart) enableLinux() error {
	// Create .desktop file in autostart directory
	desktopFile := fmt.Sprintf(`[Desktop Entry]
Type=Application
Name=%s
Exec=%s --no-gui
X-GNOME-Autostart-enabled=true
`, a.DisplayName, a.Executable)

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	dir := fmt.Sprintf("%s/.config/autostart", homeDir)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create autostart directory: %w", err)
	}

	path := fmt.Sprintf("%s/%s.desktop", dir, a.Name)
	if err := os.WriteFile(path, []byte(desktopFile), 0644); err != nil {
		return fmt.Errorf("failed to write desktop file: %w", err)
	}

	return nil
}

func (a *AutoStart) disableLinux() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	path := fmt.Sprintf("%s/.config/autostart/%s.desktop", homeDir, a.Name)
	err = os.Remove(path)
	if errors.Is(err, os.ErrNotExist) {
		// File doesn't exist, which is what we want
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to remove desktop file: %w", err)
	}

	return nil
}

// macOS implementation
func (a *AutoStart) isEnabledMacOS() bool {
	// TODO: Check if a LaunchAgent exists for this app
	return false
}

func (a *AutoStart) enableMacOS() error {
	// TODO: Create a LaunchAgent plist file
	return nil
}

func (a *AutoStart) disableMacOS() error {
	// TODO: Remove the LaunchAgent plist file
	return nil
}
