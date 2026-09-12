package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/zencodecode/authorizer-service/internal/config"
	"github.com/zencodecode/authorizer-service/internal/domain/service"
	"github.com/zencodecode/authorizer-service/internal/infrastructure/driver/database/postgres"
	appLog "github.com/zencodecode/authorizer-service/internal/infrastructure/driver/logger"
	pgRepo "github.com/zencodecode/authorizer-service/internal/infrastructure/persistence/postgres/repository"
	seedSvc "github.com/zencodecode/authorizer-service/internal/infrastructure/seed"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("seed failed: %v", err)
	}
}

func run() error {
	if err := godotenv.Load(); err != nil {
		log.Printf("warning: .env not found: %v", err)
	}

	pgCfg, err := config.LoadPostgresConfig()
	if err != nil {
		return err
	}

	seedCfg, err := config.LoadSeedConfig()
	if err != nil {
		return err
	}

	logger := appLog.New(0)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	connectCtx, cancelConnect := context.WithTimeout(ctx, 30*time.Second)
	defer cancelConnect()

	pg := postgres.NewConnection(&pgCfg, logger)
	if err := pg.Connect(connectCtx); err != nil {
		return err
	}
	defer pg.Close(context.Background(), 10*time.Second)

	logger.Info(ctx, "database connection established")
	db := pg.GetClient()
	userRepo := pgRepo.NewUserRepository(db)
	appRepo := pgRepo.NewApplicationRepository(db)
	roleRepo := pgRepo.NewRoleRepository(db)
	permRepo := pgRepo.NewPermissionRepository(db)
	rolePermRepo := pgRepo.NewRolePermissionRepository(db)
	userRoleRepo := pgRepo.NewUserRoleRepository(db)

	uc := seedSvc.NewSeederService(userRepo, appRepo, roleRepo, permRepo, rolePermRepo, userRoleRepo, logger)

	if err := uc.Seed(ctx, service.SeederParams{
		AdminEmail:      seedCfg.AdminEmail,
		AdminPassword:   seedCfg.AdminPassword,
		AdminName:       seedCfg.AdminName,
		AppClientSecret: seedCfg.AppClientSecret,
	}); err != nil {
		return err
	}

	logger.Info(ctx, "seed process completed — exiting")
	return nil
}
