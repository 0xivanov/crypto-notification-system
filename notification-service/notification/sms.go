package notification

import "log"

type SMSNotifier struct {
	apiEndpoint string
	apiKey      string
	logger      *log.Logger
}

func NewSMSNotifier(apiEndpoint, apiKey string, logger *log.Logger) *SMSNotifier {
	return &SMSNotifier{
		apiEndpoint: apiEndpoint,
		apiKey:      apiKey,
		logger:      logger,
	}
}

func (s *SMSNotifier) SendMessage(destination string, message string) error {
	// TODO
	return nil
}
