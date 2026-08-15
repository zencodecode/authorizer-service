package user

import (
	"context"

	"github.com/zencodecode/authorizer-service/internal/domain/entity"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/user"
	"github.com/zencodecode/authorizer-service/internal/domain/service"
)

type (
	CreateParams struct {
		Username string
		FullName string
		Email    string
		Phone    string
		Password string
	}
)

type createUsecase struct {
	repo   user.Repository
	logger service.Logger
}

func NewCreateUsecase(repo user.Repository, logger service.Logger) CreateUseCase {
	return &createUsecase{
		repo:   repo,
		logger: logger,
	}
}

func (uc *createUsecase) Execute(ctx context.Context, params CreateParams) (*entity.User, error) {
	return nil, nil
}
