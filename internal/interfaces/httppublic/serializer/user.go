package serializer

import (
	"github.com/google/uuid"
	"github.com/zencodecode/authorizer-service/internal/domain/entity"
)

type User struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	Name   string    `json:"name"`
}

func SerializeToUser(u entity.User) User {
	return User{
		UserID: u.ID,
		Email:  u.Email,
		Name:   u.Name,
	}
}
