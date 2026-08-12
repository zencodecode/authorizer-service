package model

import (
	"github.com/google/uuid"
	"github.com/zencodecode/authorizer-service/internal/domain/entity"
)

type ApplicationScope struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey"`
	ApplicationID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_app_scope"`
	Scope         string    `gorm:"type:varchar(255);not null;uniqueIndex:idx_app_scope"`
	Description   *string   `gorm:"type:text"`
}

func (ApplicationScope) TableName() string {
	return "application_scopes"
}

func ApplicationScopeFromEntity(e *entity.ApplicationScope) *ApplicationScope {
	return &ApplicationScope{
		ID:            e.ID,
		ApplicationID: e.ApplicationID,
		Scope:         e.Scope,
		Description:   e.Description,
	}
}

func (m *ApplicationScope) ToEntity() *entity.ApplicationScope {
	return &entity.ApplicationScope{
		ID:            m.ID,
		ApplicationID: m.ApplicationID,
		Scope:         m.Scope,
		Description:   m.Description,
	}
}
