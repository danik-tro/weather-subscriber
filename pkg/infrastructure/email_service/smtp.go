package smtp

import (
	"context"
	"fmt"
	"log"

	"gopkg.in/gomail.v2"

	domain "github.com/danik-tro/weather-subscriber/pkg/domain"
)

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

type SmtpService struct {
	config SMTPConfig
}

func NewSmtpService(config SMTPConfig) domain.EmailService {
	return &SmtpService{
		config: config,
	}
}

func (s *SmtpService) SendMessage(ctx context.Context, recipient string, subject string, body string) error {
	m := gomail.NewMessage()

	m.SetHeader("From", s.config.From)
	m.SetHeader("To", recipient)
	m.SetHeader("Subject", subject)

	m.SetBody("text/html", body)

	d := gomail.NewDialer(s.config.Host, s.config.Port, s.config.Username, s.config.Password)

	if err := d.DialAndSend(m); err != nil {
		log.Printf("Failed to send email: %v", err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
