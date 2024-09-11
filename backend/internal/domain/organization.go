package domain

import (
	"github.com/google/uuid"
)

const (
	IE  = "IE"
	LLC = "LLC"
	JSC = "JSC"
)

type Organization struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Type        *string   `json:"type,omitempty"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   *string   `json:"updated_at,omitempty"`
}

type OrganizationResponsible struct {
	ID             uuid.UUID `json:"id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	UserID         uuid.UUID `json:"user_id"`
}

type CreateOrganizationRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Type        *string `json:"type,omitempty"`
}

type AssignEmployeeToOrganizationRequest struct {
	OrganizationID uuid.UUID `json:"organization_id"`
	UserID         uuid.UUID `json:"user_id"`
}
