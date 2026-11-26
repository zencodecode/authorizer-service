package email

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/zencodecode/authorizer-service/internal/domain/service"
	"github.com/zencodecode/authorizer-service/internal/infrastructure/driver/rabbitmq"
)

type rabbitmqPublisher struct {
	ch     *amqp.Channel
	logger service.Logger
}

func NewPublisher(ch *amqp.Channel, logger service.Logger) (service.EmailSender, error) {
	if err := rabbitmq.EmailSetupTopology(ch); err != nil {
		return nil, fmt.Errorf("failed to setup topology: %w", err)
	}

	return &rabbitmqPublisher{ch: ch, logger: logger}, nil
}

func (s *rabbitmqPublisher) Send(ctx context.Context, to, subject, body string) error {
	payload := service.EmailMessage{
		To:      to,
		Subject: subject,
		Body:    body,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal email params: %w", err)
	}

	err = s.ch.PublishWithContext(ctx,
		rabbitmq.MainExchange,
		rabbitmq.MainRouting,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         data,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish email message: %w", err)
	}

	s.logger.Info(ctx, "email message published to queue",
		"to", to,
		"subject", subject,
	)

	return nil
}
