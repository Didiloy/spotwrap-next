//go:build linux
// +build linux

package autostart

import (
	"errors"
	"fmt"
	"os"
)

func isEnabled(a *AutoStart) bool {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false
	}
	path := fmt.Sprintf("%s/.config/autostart/%s.desktop", homeDir, a.Name)
	_, err = os.Stat(path)
	return err == nil
}

func enable(a *AutoStart) error {
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

func disable(a *AutoStart) error {
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
