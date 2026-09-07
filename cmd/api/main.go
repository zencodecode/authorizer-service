package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/zencodecode/authorizer-service/internal/bootstrap"
	"github.com/zencodecode/authorizer-service/internal/config"
	"github.com/zencodecode/authorizer-service/internal/infrastructure/driver/database"
	appLog "github.com/zencodecode/authorizer-service/internal/infrastructure/driver/logger"
	"github.com/zencodecode/authorizer-service/internal/infrastructure/driver/rabbitmq"
	"github.com/zencodecode/authorizer-service/internal/interfaces/event"
	"github.com/zencodecode/authorizer-service/internal/interfaces/httppublic"
	"golang.org/x/sync/errgroup"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("fatal error: %v", err)
	}
}

func run() error {
	if err := godotenv.Load(); err != nil {
		log.Printf("warning: .env not found: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	logger := appLog.New(0)
	connectCtx, cancelConnect := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelConnect()

	connManager, err := database.NewConnectionManager(connectCtx, cfg.Database, logger)
	if err != nil {
		return fmt.Errorf("initialize database connections: %w", err)
	}
	logger.Info(connectCtx, "database connections established")

	amqp := rabbitmq.NewConnection(cfg.RabbitMQ, logger)
	container := bootstrap.NewContainer(cfg, connManager, amqp, logger)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	g, gCtx := errgroup.WithContext(ctx)
	switch os.Getenv("INTERFACE") {
	case "HTTP_PUBLIC":
		container.Build(connManager.Postgres.GetClient(), connManager.Redis.GetClient())
		g.Go(func() error {
			return httppublic.Launch(gCtx, container)
		})

	case "HTTP_PRIVATE":
		// container.Build(db, redisClient)
		// httpprivate.Launch(ctx, container)

	case "EVENT":
		event.Launch(ctx, container)
	}

	serviceErr := g.Wait()
	if serviceErr != nil {
		logger.Error(gCtx, "service exited with error", "error", serviceErr.Error())
	}

	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 15*time.Second)
	logger.Info(shutdownCtx, "shutting down database connections...")
	defer cancelShutdown()

	if shutdownErr := connManager.Shutdown(shutdownCtx); shutdownErr != nil {
		logger.Error(shutdownCtx, "database shutdown error", "error", shutdownErr.Error())
		if serviceErr == nil {
			return shutdownErr
		}
	}

	return serviceErr
}
