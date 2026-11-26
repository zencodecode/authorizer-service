package email

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/zencodecode/authorizer-service/internal/domain/service"
	"github.com/zencodecode/authorizer-service/internal/infrastructure/driver/rabbitmq"
)

const maxRetries = 3

type RabbitmqConsumer struct {
	ch     *amqp.Channel
	mailer service.EmailSender
	logger service.Logger
}

func NewConsumer(ch *amqp.Channel, mailer service.EmailSender, logger service.Logger) *RabbitmqConsumer {
	return &RabbitmqConsumer{
		ch:     ch,
		mailer: mailer,
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

		case delivery, ok := <-msgs:
			if !ok {
				return fmt.Errorf("channel closed")
			}
			c.handleMessage(ctx, delivery)
		}
	}
}

func (c *RabbitmqConsumer) handleMessage(ctx context.Context, delivery amqp.Delivery) {
	defer func() {
		if r := recover(); r != nil {
			c.logger.Error(ctx, "panic while handling email message", "recover", fmt.Sprintf("%v", r))
			c.sendToFailedQueue(ctx, delivery, "panic recovered")
			delivery.Ack(false)
		}
	}()

	var msg service.EmailMessage
	if err := json.Unmarshal(delivery.Body, &msg); err != nil {
		c.logger.Error(ctx, "failed to unmarshal email message",
			"error", err.Error(),
		)
		c.sendToFailedQueue(ctx, delivery, "unmarshal error: "+err.Error())
		delivery.Ack(false)
		return
	}

	if err := c.mailer.Send(ctx, msg.To, msg.Subject, msg.Body); err != nil {
		retryCount := getRetryCount(delivery)

		if retryCount >= maxRetries {
			c.logger.Error(ctx, "email send failed permanently, giving up",
				"to", msg.To,
				"retries", retryCount,
				"error", err.Error(),
			)
			c.sendToFailedQueue(ctx, delivery, err.Error())
			delivery.Ack(false)
			return
		}

		c.logger.Error(ctx, "failed to send email, will retry",
			"to", msg.To,
			"retries", retryCount+1,
			"error", err.Error(),
		)

		delivery.Nack(false, false)
		return
	}

	delivery.Ack(false)
	c.logger.Info(ctx, "email sent successfully", "to", msg.To)
}

func (c *RabbitmqConsumer) sendToFailedQueue(ctx context.Context, delivery amqp.Delivery, reason string) {
	headers := amqp.Table{}
	for k, v := range delivery.Headers {
		headers[k] = v
	}
	headers["x-failure-reason"] = reason

	err := c.ch.PublishWithContext(ctx,
		rabbitmq.FailedExchange,
		rabbitmq.FailedRouting,
		false,
		false,
		amqp.Publishing{
			ContentType:  delivery.ContentType,
			Body:         delivery.Body,
			Headers:      headers,
			DeliveryMode: amqp.Persistent,
		},
	)
	if err != nil {
		c.logger.Error(ctx, "failed to publish message to failed queue", "error", err.Error())
	}
}

func getRetryCount(delivery amqp.Delivery) int64 {
	deaths, ok := delivery.Headers["x-death"].([]interface{})
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
