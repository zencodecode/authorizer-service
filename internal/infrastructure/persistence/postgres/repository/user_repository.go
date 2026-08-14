package repository

import (
	"context"
	"errors"

	"github.com/zencodecode/authorizer-service/internal/domain/entity"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/user"
	"github.com/zencodecode/authorizer-service/internal/infrastructure/persistence/postgres/model"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) user.Repository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	userModel := model.UserFromEntity(user)
	return r.db.WithContext(ctx).Create(userModel).Error
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*entity.User, error) {
	var userModel model.User
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&userModel)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, user.ErrNotFound
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return userModel.ToEntity(), nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	var userModel model.User
	result := r.db.WithContext(ctx).Where("email = ?", email).First(&userModel)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, user.ErrNotFound
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return userModel.ToEntity(), nil
}

func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
	userModel := model.UserFromEntity(user)
	return r.db.WithContext(ctx).Save(userModel).Error
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.User{}).Error
}

func (r *userRepository) List(ctx context.Context, limit, offset int) ([]*entity.User, error) {
	var usersModel []model.User
	result := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&usersModel)
	if result.Error != nil {
		return nil, result.Error
	}

	usersEntity := make([]*entity.User, len(usersModel))
	for i, m := range usersModel {
		usersEntity[i] = m.ToEntity()
	}
	return usersEntity, nil
}
