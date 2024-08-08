package model

type User struct {
	UserID              string              `bson:"userID"`
	Tickers             []string            `bson:"tickers"`
	NotificationOptions NotificationOptions `bson:"notificationOptions"`
}

type NotificationOptions struct {
	SlackWebhookURL string `bson:"slackWebhookURL"`
	Email           string `bson:"email"`
	PhoneNumber     string `bson:"phoneNumber"`
}
