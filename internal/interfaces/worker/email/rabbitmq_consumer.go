package email

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/zencodecode/authorizer-service/internal/domain/service"
	"github.com/zencodecode/authorizer-service/internal/infrastructure/driver/email"
)

type RabbitmqConsumer struct {
	ch     *amqp.Channel
	smtp   *email.SMTPSender
	logger service.Logger
}

func NewConsumer(ch *amqp.Channel, smtp *email.SMTPSender, logger service.Logger) *RabbitmqConsumer {
	return &RabbitmqConsumer{
		ch:     ch,
		smtp:   smtp,
		logger: logger,
	}
}

func (c *RabbitmqConsumer) Start(ctx context.Context) error {
	msgs, err := c.ch.Consume(
		"email.send",
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to start consuming: %w", err)
	}

	c.logger.Info(ctx, "email consumer started, waiting for messages...")

	for {
		select {
		case <-ctx.Done():
			c.logger.Info(ctx, "email consumer stopped")
			return nil

		case msg, ok := <-msgs:
			if !ok {
				return fmt.Errorf("channel closed")
			}
			c.handleMessage(ctx, msg)
		}
	}
}

func (c *RabbitmqConsumer) handleMessage(ctx context.Context, msg amqp.Delivery) {
	var params service.SendEmailParams
	if err := json.Unmarshal(msg.Body, &params); err != nil {
		c.logger.Error(ctx, "failed to unmarshal email message",
			"error", err.Error(),
		)
		msg.Nack(false, false)
		return
	}
	if err := c.smtp.Send(ctx, params.To, params.Subject, params.Body); err != nil {
		c.logger.Error(ctx, "failed to send email",
			"to", params.To,
			"error", err.Error(),
		)
		msg.Nack(false, true)
		return
	}

	msg.Ack(false)
	c.logger.Info(ctx, "email sent successfully", "to", params.To)
}
