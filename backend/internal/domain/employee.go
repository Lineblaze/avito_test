package domain

import (
	"github.com/google/uuid"
)

type Employee struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	FirstName *string   `json:"first_name"`
	LastName  *string   `json:"last_name"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt *string   `json:"updated_at,omitempty"`
}

type CreateEmployeeRequest struct {
	Username  string  `json:"username"`
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
}
