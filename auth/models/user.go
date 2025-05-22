package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	// gorm.Model
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Email    string    `gorm:"unique;not null"`
	Password string    `gorm:"not null"`
	Tier     string    `gorm:"default:'basic'"`
}

type SessionMetadata struct {
	UserID    uuid.UUID `json:"user_id"`
	Email     string    `json:"email"`
	Tier      string    `json:"tier"`
	LastLogin time.Time `json:"last_login"`
	IPAddress string    `json:"ip_address"`
}
