package serializer

import (
	"github.com/google/uuid"
	"github.com/zencodecode/authorizer-service/internal/domain/entity"
)

type Organization struct {
	OrganizationID uuid.UUID `json:"organization_id"`
	Name           string    `json:"name"`
	Slug           string    `json:"slug"`
}

func SerializeToOrganization(o entity.Organization) Organization {
	return Organization{
		OrganizationID: o.ID,
		Name:           o.Name,
		Slug:           o.Slug,
	}
}
