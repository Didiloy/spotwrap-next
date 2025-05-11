package notifications

import (
	_ "embed"
	"os"
	"path/filepath"
	"sync"
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

// Notify sends a system notification with the given title and message
func Notify(title, message string) error {
	return notify(title, message)
}
