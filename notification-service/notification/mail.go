package notification

import (
	"log"

	"github.com/0xivanov/crypto-notification-system/common/model"
	"github.com/wneessen/go-mail"
)

type MailNotifier struct {
	client *mail.Client
	from   string
	logger *log.Logger
}

func NewMailNotifier(host, username, password, from string, logger *log.Logger) *MailNotifier {
	c, err := mail.NewClient(
		host,
		mail.WithPort(25),
		mail.WithPort(25), mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(username),
		mail.WithPassword(password),
		mail.WithTLSPolicy(mail.TLSMandatory))
	if err != nil {
		logger.Fatalf("[ERROR] Failed to create mail client: %v", err)
	}

	return &MailNotifier{
		client: c,
		from:   from,
		logger: logger,
	}
}

func (n *MailNotifier) SendNotification(message string, options model.NotificationOptions) error {
	if options.Email == "" {
		n.logger.Println("[INFO] Email is empty")
		return nil
	}
	m := mail.NewMsg()
	if err := m.From(n.from); err != nil {
		n.logger.Printf("[ERROR] Failed to set From address: %s", err)
		return err
	}
	if err := m.To(options.Email); err != nil {
		n.logger.Printf("[ERROR] Failed to set To address: %s", err)
		return err
	}

	m.Subject("Crypto updates")
	m.SetBodyString(mail.TypeTextPlain, message)
	if err := n.client.DialAndSend(m); err != nil {
		n.logger.Fatalf("[ERROR] Failed to send mail: %s", err)
	}
	return nil
}
