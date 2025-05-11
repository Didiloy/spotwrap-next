//go:build windows
// +build windows

package notifications

import (
	"github.com/go-toast/toast"
)

func notify(title, message string) error {
	notification := toast.Notification{
		AppID:   appName,
		Title:   title,
		Message: message,
		Icon:    getIconPath(),
	}
	return notification.Push()
}
