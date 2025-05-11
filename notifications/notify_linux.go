//go:build linux
// +build linux

package notifications

import (
	"os/exec"
)

func notify(title, message string) error {
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
