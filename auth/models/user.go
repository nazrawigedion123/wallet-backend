package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Tier     string `gorm:"default:'basic'"`
}

type SessionMetadata struct {
	UserID    uint      `json:"user_id"`
	Email     string    `json:"email"`
	Tier      string    `json:"tier"`
	LastLogin time.Time `json:"last_login"`
	IPAddress string    `json:"ip_address"`
}
