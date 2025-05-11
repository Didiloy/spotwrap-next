//go:build !windows && !linux
// +build !windows,!linux

package notifications

import (
	"os/exec"
)

func notify(title, message string) error {
	// Use default notification system for other platforms
	cmd := exec.Command("notify-send",
		"-a", appName,
		"-i", getIconPath(),
		title,
		message)
	return cmd.Run()
}
