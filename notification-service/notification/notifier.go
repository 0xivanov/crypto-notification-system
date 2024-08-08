package notification

import "github.com/0xivanov/crypto-notification-system/notification-service/model"

type Notifier interface {
	SendNotification(message string, options model.NotificationOptions) error
}
