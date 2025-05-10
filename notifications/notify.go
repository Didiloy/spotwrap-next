package notifications

import (
	_ "embed"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/go-toast/toast"
)

const appName = "Spotwrap Next"

//go:embed assets/appicon.png
var appIconData []byte

var (
	iconPath string
	iconOnce sync.Once
)

func getIconPath() string {
	iconOnce.Do(func() {
		dir := os.TempDir()
		path := filepath.Join(dir, "spotwrap-appicon.png")
		_ = os.WriteFile(path, appIconData, 0644)
		iconPath = path
	})
	return iconPath
}

func Notify(title, message string) error {
	switch runtime.GOOS {
	case "linux":
		return notifyLinux(title, message)
	case "windows":
		return notifyWindows(title, message)
	default:
		// Use default notification system for other platforms
		return notifyDefault(title, message)
	}
}

func notifyLinux(title, message string) error {
	// Use notify-send with:
	// -t 0: make notifications persistent
	// -a: specify application name
	// -i: specify icon
	cmd := exec.Command("notify-send",
		"-t", "0",
		"-a", appName,
		"-i", getIconPath(),
		title,
		message)
	return cmd.Run()
}

func notifyDefault(title, message string) error {
	// Use default notification system for other platforms
	cmd := exec.Command("notify-send",
		"-a", appName,
		"-i", getIconPath(),
		title,
		message)
	return cmd.Run()
}

func notifyWindows(title, message string) error {
	notification := toast.Notification{
		AppID:   appName,
		Title:   title,
		Message: message,
		Icon:    getIconPath(),
	}
	return notification.Push()
}
