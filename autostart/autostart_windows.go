//go:build windows
// +build windows

package autostart

import (
	"fmt"

	"golang.org/x/sys/windows/registry"
)

func isEnabled(a *AutoStart) bool {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.QUERY_VALUE)
	if err != nil {
		return false
	}
	defer k.Close()

	val, _, err := k.GetStringValue(a.DisplayName)
	if err != nil {
		return false
	}

	return val == fmt.Sprintf(`"%s" --no-gui`, a.Executable)
}

func enable(a *AutoStart) error {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("failed to open registry key: %w", err)
	}
	defer k.Close()

	err = k.SetStringValue(a.DisplayName, fmt.Sprintf(`"%s" --no-gui`, a.Executable))
	if err != nil {
		return fmt.Errorf("failed to set registry value: %w", err)
	}

	return nil
}

func disable(a *AutoStart) error {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("failed to open registry key: %w", err)
	}
	defer k.Close()

	err = k.DeleteValue(a.DisplayName)
	if err != nil && err != registry.ErrNotExist {
		return fmt.Errorf("failed to delete registry value: %w", err)
	}

	return nil
}
