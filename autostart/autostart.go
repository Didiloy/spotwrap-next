// Package autostart provides functionality to configure the application to start automatically on system boot
package autostart

import (
	"errors"
	"os"
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
	return isEnabled(a)
}

// Enable configures the application to start automatically on system boot
func (a *AutoStart) Enable() error {
	return enable(a)
}

// Disable removes the configuration for the application to start automatically
func (a *AutoStart) Disable() error {
	return disable(a)
}
