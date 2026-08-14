package bootstrap

import (
	"github.com/zencodecode/authorizer-service/internal/config"
	"github.com/zencodecode/authorizer-service/internal/domain/service"
	"github.com/zencodecode/authorizer-service/internal/infrastructure/driver/database"
)

type Container struct {
	Config      config.Config
	ConnManager *database.ConnectionManager
	Logger      service.Logger
}

func NewContainer(cfg config.Config, cm *database.ConnectionManager, logger service.Logger) *Container {
	return &Container{
		Config:      cfg,
		ConnManager: cm,
		Logger:      logger,
	}
}
