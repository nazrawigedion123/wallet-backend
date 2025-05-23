// models/webhook_payload.go

package models

import (
	"time"

	"github.com/nazrawigedion123/wallet-backend/webhook/utils"
	"gorm.io/gorm"
)

type WebhookEvent struct {
	ID        uint        `gorm:"primaryKey"`
	EventID   string      `gorm:"uniqueIndex;not null"` // For idempotency
	Type      string      `gorm:"not null"`             // bill_payment, wallet_credit, etc.
	UserID    string      `gorm:"index;not null"`
	Amount    float64     `gorm:"not null"`
	Timestamp time.Time   `gorm:"not null"`
	Metadata  utils.JSONB `gorm:"type:jsonb"`        // Custom type for map[string]string
	Status    string      `gorm:"default:'pending'"` // processed, failed, etc.
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
type Metadata struct {
	TransactionID uint
}

type IncomingWebhook struct {
	EventID   string    `json:"event_id" validate:"required"`
	Type      string    `json:"type" validate:"required"`
	UserID    string    `json:"user_id" validate:"required"`
	Amount    float64   `json:"amount" validate:"required"`
	Timestamp time.Time `json:"timestamp"`
	Metadata  Metadata  `json:"metadata"`
	Status    string    `json:"status" validate:"required"`
}

// type Metadata struct {
// 	TransactionID uint `json:"transaction_id"`
// }

// type IncomingWebhook struct {
// 	EventID   string    `json:"reference" validate:"required"`  // Changed to match JSON
// 	Type      string    `json:"event" validate:"required"`      // Changed to match JSON
// 	UserID    string    `json:"user_id" validate:"required"`    // Not in JSON sample!
// 	Amount    float64   `json:"amount" validate:"required"`
// 	Timestamp time.Time `json:"timestamp"`                      // Not in JSON sample!
// 	Currency  string    `json:"currency"`                       // Added new field
// 	Status    string    `json:"status" validate:"required"`     // Added new field
// 	Metadata  Metadata  `json:"metadata"`
// }
