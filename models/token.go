package models

import (
	"time"

	"github.com/google/uuid"
)

type Token struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	UserID    uuid.UUID
	ExpiresAt time.Time
}
