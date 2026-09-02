package email

import (
	"context"
	"fmt"
	"net/smtp"
	"strings"

	"github.com/zencodecode/authorizer-service/internal/config"
	"github.com/zencodecode/authorizer-service/internal/domain/service"
)

type smtpSender struct {
	cfg config.SMTP
}

func New(cfg config.SMTP) service.EmailSender {
	return &smtpSender{cfg: cfg}
}

func (s *smtpSender) Send(ctx context.Context, params service.SendEmailParams) error {
	addr := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)

	from := fmt.Sprintf("%s <%s>", s.cfg.FromName, s.cfg.From)

	headers := []string{
		fmt.Sprintf("From: %s", from),
		fmt.Sprintf("To: %s", params.To),
		fmt.Sprintf("Subject: %s", params.Subject),
		"MIME-Version: 1.0",
		"Content-Type: text/html; charset=\"UTF-8\"",
	}

	msg := []byte(strings.Join(headers, "\r\n") + "\r\n\r\n" + params.Body)

	var auth smtp.Auth
	if s.cfg.Username != "" {
		auth = smtp.PlainAuth("", s.cfg.Username, s.cfg.Password, s.cfg.Host)
	}

	return smtp.SendMail(addr, auth, s.cfg.From, []string{params.To}, msg)
}
