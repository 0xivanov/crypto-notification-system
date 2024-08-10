package model

type User struct {
	UserID              string              `bson:"userID" json:"userID"`
	Tickers             []TickerSettings    `bson:"tickers" json:"tickers"`
	NotificationOptions NotificationOptions `bson:"notificationOptions" json:"notificationOptions"`
}

type TickerSettings struct {
	Symbol          string  `bson:"symbol" json:"symbol"`
	ChangeThreshold float64 `bson:"changeThreshold" json:"changeThreshold"`
}

type NotificationOptions struct {
	SlackWebhookURL string `bson:"slackWebhookURL" json:"slackWebhookURL"`
	Email           string `bson:"email" json:"email"`
	PhoneNumber     string `bson:"phoneNumber" json:"phoneNumber"`
}
