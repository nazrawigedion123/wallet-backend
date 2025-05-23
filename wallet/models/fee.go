package models

import (
	"time"
)

type UserTier string

const (
	BasicTier   UserTier = "basic"
	PremiumTier UserTier = "premium"
)

type FeeConfig struct {
	TransactionType string    `json:"transaction_type"`
	Tier            UserTier  `json:"tier"`
	BasePercent     float64   `json:"base_percent"`
	Cap             float64   `json:"cap"`
	Floor           float64   `json:"floor"`
	PeakStart       time.Time `json:"peak_start"`
	PeakEnd         time.Time `json:"peak_end"`
	PeakSurcharge   float64   `json:"peak_surcharge"` // extra percent during peak
}
