package service

import "context"

type EmailMessage struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type EmailSender interface {
	Send(ctx context.Context, to, subject, body string) error
}
