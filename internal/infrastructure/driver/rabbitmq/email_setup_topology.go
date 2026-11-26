package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	MainExchange = "authorizer.email"
	MainQueue    = "email.send"
	MainRouting  = "email.send"

	RetryExchange = "email.retry.exchange"
	RetryQueue    = "email.send.retry"
	RetryRouting  = "email.send.retry"

	FailedExchange = "email.failed.exchange"
	FailedQueue    = "email.send.failed"
	FailedRouting  = "email.send.failed"

	RetryTTLMs = 30000
)

func EmailSetupTopology(ch *amqp.Channel) error {
	if err := ch.ExchangeDeclare(
		MainExchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return err
	}

	if err := ch.ExchangeDeclare(
		RetryExchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return err
	}

	_, err := ch.QueueDeclare(
		RetryQueue,
		true,
		false,
		false,
		false,
		amqp.Table{
			"x-dead-letter-exchange":    "",
			"x-dead-letter-routing-key": MainQueue,
			"x-message-ttl":             int32(RetryTTLMs),
		},
	)
	if err != nil {
		return err
	}

	if err := ch.QueueBind(
		RetryQueue,
		RetryRouting,
		RetryExchange,
		false,
		nil,
	); err != nil {
		return err
	}

	if err := ch.ExchangeDeclare(
		FailedExchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return err
	}

	if _, err := ch.QueueDeclare(
		FailedQueue,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return err
	}

	if err := ch.QueueBind(
		FailedQueue,
		FailedRouting,
		FailedExchange,
		false,
		nil,
	); err != nil {
		return err
	}

	if _, err := ch.QueueDeclare(
		MainQueue,
		true,
		false,
		false,
		false,
		amqp.Table{
			"x-dead-letter-exchange":    RetryExchange,
			"x-dead-letter-routing-key": RetryQueue,
		},
	); err != nil {
		return err
	}

	if err := ch.QueueBind(
		MainQueue,
		MainRouting,
		MainExchange,
		false,
		nil,
	); err != nil {
		return err
	}

	return nil
}
