package notification

import (
	"log"

	"github.com/0xivanov/crypto-notification-system/notification-service/model"
	"github.com/slack-go/slack"
)

type SlackNotifier struct {
	logger *log.Logger
}

func NewSlackNotifier(logger *log.Logger) *SlackNotifier {
	return &SlackNotifier{logger: logger}
}

func (s *SlackNotifier) SendNotification(message string, options model.NotificationOptions) error {
	if options.SlackWebhookURL == "" {
		s.logger.Println("[INFO] Slack webhook URL is empty")
		return nil
	}

	payload := slack.WebhookMessage{
		Username: "Crypto Notification System",
		Text:     message,
	}

	return slack.PostWebhook(options.SlackWebhookURL, &payload)
}
