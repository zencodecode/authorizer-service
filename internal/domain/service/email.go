package service

import "context"

type SendEmailParams struct {
	To      string
	Subject string
	Body    string
}

type EmailSender interface {
	Send(ctx context.Context, params SendEmailParams) error
}
