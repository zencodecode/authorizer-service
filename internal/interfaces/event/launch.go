package event

import (
	"context"
	"fmt"
	"sync"

	"github.com/zencodecode/authorizer-service/internal/bootstrap"
	"github.com/zencodecode/authorizer-service/internal/infrastructure/driver/email"
	emailevent "github.com/zencodecode/authorizer-service/internal/interfaces/event/email"
)

// Launch starts all event consumers.
// This is a blocking call — it runs until ctx is cancelled or a consumer fails.
// Does not require a port. Connects to RabbitMQ and waits for messages.
func Launch(ctx context.Context, c *bootstrap.Container) error {
	c.Logger.Info(ctx, "starting event workers")

	// ──────── Email Consumer ────────
	consumerCh, err := c.Amqp.NewChannel()
	if err != nil {
		return fmt.Errorf("failed to create consumer channel: %w", err)
	}

	smtpSender := email.NewSMTPSender(c.Config.SMTP)
	emailConsumer := emailevent.NewConsumer(consumerCh, smtpSender, c.Logger)

	// ──────── Start all consumers ────────
	var wg sync.WaitGroup
	errCh := make(chan error, 1)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := emailConsumer.Start(ctx); err != nil {
			c.Logger.Error(ctx, "email consumer error", "error", err.Error())
			errCh <- fmt.Errorf("email consumer: %w", err)
		}
	}()

	// Nanti tambah consumer lain di sini:
	// wg.Add(1)
	// go func() {
	//     defer wg.Done()
	//     notifConsumer := notification.NewConsumer(...)
	//     if err := notifConsumer.Start(ctx); err != nil { ... }
	// }()

	c.Logger.Info(ctx, "all event workers started")

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		c.Logger.Info(ctx, "shutting down event workers")
		wg.Wait()
		return nil
	}
}
