package notification

import "github.com/0xivanov/crypto-notification-system/common/model"

type Notifier interface {
	SendNotification(message string, options model.NotificationOptions) error
}
