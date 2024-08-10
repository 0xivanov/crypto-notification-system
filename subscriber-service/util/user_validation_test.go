package util

import "testing"

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		email    string
		expected bool
	}{
		{"test@example.com", true},
		{"user.name+tag+sorting@example.com", false},
		{"user.name@example.co.uk", true},
		{"user.name@.com", false},
		{"plainaddress", false},
		{"@missingusername.com", false},
		{"user@.com", false},
		{"user@com", false},
		{"user@domain.com.", false},
	}

	for _, test := range tests {
		result := IsValidEmail(test.email)
		if result != test.expected {
			t.Errorf("IsValidEmail(%s) = %t; expected %t", test.email, result, test.expected)
		}
	}
}

func TestIsValidPhoneNumber(t *testing.T) {
	tests := []struct {
		phoneNumber string
		expected    bool
	}{
		{"+1234567890", true},
		{"+19876543210", true},
		{"1234567890", true},
		{"+1", false},
		{"", false},
		{"+1234567890123456", false},
		{"abc1234567", false},
		{"+359882725331", true},
	}

	for _, test := range tests {
		result := IsValidPhoneNumber(test.phoneNumber)
		if result != test.expected {
			t.Errorf("IsValidPhoneNumber(%s) = %t; expected %t", test.phoneNumber, result, test.expected)
		}
	}
}

func TestIsValidSlackWebhook(t *testing.T) {
	tests := []struct {
		webhookUrl string
		expected   bool
	}{
		{"https://hooks.slack.com/services/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", true},
		{"https://hooks.slack.com/services/T07FQUNKWQ3/B07GDPJAP08/", false},
		{"https://hooks.slack.com/services/invalid/url/here", false},
		{"http://hooks.slack.com/services/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", false}, // Not https
	}

	for _, test := range tests {
		result := IsValidSlackWebhook(test.webhookUrl)
		if result != test.expected {
			t.Errorf("IsValidSlackWebhook(%s) = %t; expected %t", test.webhookUrl, result, test.expected)
		}
	}
}
