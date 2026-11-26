package email

import (
	"context"
	"fmt"
	"net/smtp"
	"strings"

	"github.com/zencodecode/authorizer-service/internal/config"
)

type SMTPSender struct {
	cfg config.SMTP
}

func NewSMTPSender(cfg config.SMTP) *SMTPSender {
	return &SMTPSender{cfg: cfg}
}

func (s *SMTPSender) Send(ctx context.Context, to, subject, body string) error {
	addr := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)

	from := fmt.Sprintf("%s <%s>", s.cfg.FromName, s.cfg.From)

	headers := []string{
		fmt.Sprintf("From: %s", from),
		fmt.Sprintf("To: %s", to),
		fmt.Sprintf("Subject: %s", subject),
		"MIME-Version: 1.0",
		"Content-Type: text/html; charset=\"UTF-8\"",
	}

	msg := []byte(strings.Join(headers, "\r\n") + "\r\n\r\n" + body)

	var auth smtp.Auth
	if s.cfg.Username != "" {
		auth = smtp.PlainAuth("", s.cfg.Username, s.cfg.Password, s.cfg.Host)
	}

	return smtp.SendMail(addr, auth, s.cfg.From, []string{to}, msg)
}
