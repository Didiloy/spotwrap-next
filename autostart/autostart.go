package autostart

import (
	"os"
	"runtime"
)

type AutoStart struct {
	Name        string
	DisplayName string
	Executable  string
}

func New(name, displayName string) *AutoStart {
	executable, _ := os.Executable()
	return &AutoStart{
		Name:        name,
		DisplayName: displayName,
		Executable:  executable,
	}
}

func (a *AutoStart) IsEnabled() bool {
	switch runtime.GOOS {
	case "windows":
		return a.isEnabledWindows()
	case "linux":
		return a.isEnabledLinux()
	default:
		return false
	}
}

func (a *AutoStart) Enable() error {
	switch runtime.GOOS {
	case "windows":
		return a.enableWindows()
	case "linux":
		return a.enableLinux()
	default:
		return nil
	}
}

func (a *AutoStart) Disable() error {
	switch runtime.GOOS {
	case "windows":
		return a.disableWindows()
	case "linux":
		return a.disableLinux()
	default:
		return nil
	}
}

// Windows implementation
func (a *AutoStart) isEnabledWindows() bool {
	// Check registry for startup entry
	return false
}

func (a *AutoStart) enableWindows() error {
	// Create registry entry
	// Example: HKCU\Software\Microsoft\Windows\CurrentVersion\Run
	return nil
}

func (a *AutoStart) disableWindows() error {
	// Remove registry entry
	return nil
}

// Linux implementation
func (a *AutoStart) isEnabledLinux() bool {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false
	}
	path := homeDir + "/.config/autostart/" + a.Name + ".desktop"
	_, err = os.Stat(path)
	return err == nil
}

func (a *AutoStart) enableLinux() error {
	// Create .desktop file in autostart directory
	desktopFile := `[Desktop Entry]
Type=Application
Name=` + a.DisplayName + `
Exec=` + a.Executable + ` --no-gui
X-GNOME-Autostart-enabled=true
`
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	dir := homeDir + "/.config/autostart"
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	path := dir + "/" + a.Name + ".desktop"
	return os.WriteFile(path, []byte(desktopFile), 0644) // Write to ~/.config/autostart/appname.desktop
}

func (a *AutoStart) disableLinux() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	path := homeDir + "/.config/autostart/" + a.Name + ".desktop"
	return os.Remove(path)
}
