package util

import "regexp"

func IsValidEmail(email string) bool {
	var emailRegex = regexp.MustCompile(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`)

	return emailRegex.MatchString(email)
}

func IsValidPhoneNumber(phoneNumber string) bool {
	var phoneNumberRegex = regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)

	return phoneNumberRegex.MatchString(phoneNumber)
}

func IsValidSlackWebhook(webhookUrl string) bool {
	var slackWebhookRegex = regexp.MustCompile(`^https:\/\/hooks\.slack\.com\/services\/[A-Z0-9]+\/[A-Z0-9]+\/[a-zA-Z0-9]+$`)

	return slackWebhookRegex.MatchString(webhookUrl)
}
