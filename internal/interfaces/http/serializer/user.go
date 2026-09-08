package serializer

import (
	"github.com/google/uuid"
	"github.com/zencodecode/authorizer-service/internal/domain/entity"
)

type User struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
	Name  string    `json:"name"`
}

func SerializeToUser(u entity.User) User {
	return User{
		ID:    u.ID,
		Email: u.Email,
		Name:  u.Name,
	}
}
