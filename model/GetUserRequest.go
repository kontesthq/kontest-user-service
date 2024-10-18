package model

import "github.com/google/uuid"

type GetUserRequest struct {
	UserID uuid.UUID
}
