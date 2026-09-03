package email

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/zencodecode/authorizer-service/internal/domain/service"
	"github.com/zencodecode/authorizer-service/internal/infrastructure/driver/email"
	"github.com/zencodecode/authorizer-service/internal/infrastructure/driver/rabbitmq"
)

const maxRetries = 3

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
	if err := rabbitmq.EmailSetupTopology(c.ch); err != nil {
		return fmt.Errorf("failed to setup topology: %w", err)
	}

	if err := c.ch.Qos(1, 0, false); err != nil {
		return fmt.Errorf("failed to set qos: %w", err)
	}

	msgs, err := c.ch.Consume(
		rabbitmq.MainQueue,
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
	defer func() {
		if r := recover(); r != nil {
			c.logger.Error(ctx, "panic while handling email message", "recover", fmt.Sprintf("%v", r))
			c.sendToFailedQueue(ctx, msg, "panic recovered")
		}
	}()

	var params service.SendEmailParams
	if err := json.Unmarshal(msg.Body, &params); err != nil {
		c.logger.Error(ctx, "failed to unmarshal email message",
			"error", err.Error(),
		)
		c.sendToFailedQueue(ctx, msg, "unmarshal error: "+err.Error())
		msg.Ack(false)
		return
	}

	if err := c.smtp.Send(ctx, params.To, params.Subject, params.Body); err != nil {
		retryCount := getRetryCount(msg)

		if retryCount >= maxRetries {
			c.logger.Error(ctx, "email send failed permanently, giving up",
				"to", params.To,
				"retries", retryCount,
				"error", err.Error(),
			)
			c.sendToFailedQueue(ctx, msg, err.Error())
			msg.Ack(false)
			return
		}

		c.logger.Error(ctx, "failed to send email, will retry",
			"to", params.To,
			"retries", retryCount+1,
			"error", err.Error(),
		)

		msg.Nack(false, false)
		return
	}

	msg.Ack(false)
	c.logger.Info(ctx, "email sent successfully", "to", params.To)
}

func (c *RabbitmqConsumer) sendToFailedQueue(ctx context.Context, msg amqp.Delivery, reason string) {
	headers := amqp.Table{}
	for k, v := range msg.Headers {
		headers[k] = v
	}
	headers["x-failure-reason"] = reason

	err := c.ch.PublishWithContext(ctx,
		rabbitmq.FailedExchange,
		rabbitmq.FailedQueue,
		false,
		false,
		amqp.Publishing{
			ContentType:  msg.ContentType,
			Body:         msg.Body,
			Headers:      headers,
			DeliveryMode: amqp.Persistent,
		},
	)
	if err != nil {
		c.logger.Error(ctx, "failed to publish message to failed queue", "error", err.Error())
	}
}

func getRetryCount(msg amqp.Delivery) int64 {
	deaths, ok := msg.Headers["x-death"].([]interface{})
	if !ok || len(deaths) == 0 {
		return 0
	}
	death, ok := deaths[0].(amqp.Table)
	if !ok {
		return 0
	}
	count, _ := death["count"].(int64)
	return count
}
