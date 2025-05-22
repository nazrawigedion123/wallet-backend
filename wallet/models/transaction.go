package models

import (
	"github.com/google/uuid"
	"github.com/nazrawigedion123/wallet-backend/auth/models"
	"gorm.io/gorm"
)

type TransactionType string
type TransactionStatus string

const (
	DepositTransaction  TransactionType = "deposit"
	WithdrawTransaction TransactionType = "withdraw"

	StatusPending TransactionStatus = "pending"
	StatusSuccess TransactionStatus = "success"
	StatusFailed  TransactionStatus = "failed"
)

type Transaction struct {
	gorm.Model
	UserID uuid.UUID         `json:"user_id" gorm:"type:uuid;not null;index"`
	Amount float64           `json:"amount" gorm:"not null"`
	Type   TransactionType   `json:"type" gorm:"type:varchar(20);not null"`
	Status TransactionStatus `json:"status" gorm:"type:varchar(20);not null;default:'pending'"`

	User models.User `json:"-" gorm:"foreignKey:UserID;references:ID"`
}

type WalletBalance struct {
	UserID  uuid.UUID `json:"user_id" gorm:"type:uuid;primaryKey"`
	Balance float64   `json:"balance" gorm:"not null;default:0"`

	User models.User `gorm:"foreignKey:UserID;references:ID"`
}
