package notifications

import (
	"github.com/gen2brain/beeep"
)

func Notify(title, message string) error {
	err := beeep.Notify(title, message, "")
	return err

	// Special cases for different platforms
	// switch runtime.GOOS {
	// case "linux":
	// 	// Additional Linux-specific notification handling if needed
	// 	return notifyLinux(title, message)
	// case "windows":
	// 	// Additional Windows-specific notification handling if needed
	// 	return notifyWindows(title, message)
	// default:
	// 	return nil
	// }
}

// func notifyLinux(title, message string) error {
// 	return nil
// }

// func notifyWindows(title, message string) error {
// 	return nil
// }
