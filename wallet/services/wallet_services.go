package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/nazrawigedion123/wallet-backend/wallet/models"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WalletService struct {
	redisClient *redis.Client
	dbWriter    chan models.Transaction
	db          *gorm.DB
}

var ctx = context.Background()

func NewWalletService(db *gorm.DB, redisClient *redis.Client) *WalletService {
	service := &WalletService{
		db:          db,
		redisClient: redisClient,

		dbWriter: make(chan models.Transaction, 100), // Buffered to avoid blocking
	}

	// Start a goroutine to consume the channel
	go func() {
		for txn := range service.dbWriter {

			if err := service.db.Create(&txn).Error; err != nil {
				log.Printf("failed to save transaction: %v", err)
			} else {
				fmt.Println("Transaction saved:", txn.ID)

				// Now preload the user if needed
				if err := service.db.Preload("User").First(&txn, txn.ID).Error; err != nil {
					log.Printf("failed to preload user: %v", err)
				} else {
					fmt.Printf("User loaded: %v\n", txn.User)
				}
			}
		}
	}()

	return service
}

func (ws *WalletService) GetBalance(userID uuid.UUID) (float64, error) {
	val, err := ws.redisClient.Get(ctx, ws.balanceKey(userID)).Result()
	if err == redis.Nil {
		// Redis miss: fallback to DB
		var wb models.WalletBalance
		if dbErr := ws.db.First(&wb, "user_id = ?", userID).Error; dbErr != nil {
			return 0.0, dbErr
		}
		// Optionally repopulate Redis
		_ = ws.redisClient.Set(ctx, ws.balanceKey(userID), fmt.Sprintf("%f", wb.Balance), 0).Err()
		return wb.Balance, nil
	} else if err != nil {
		return 0.0, err
	}

	var balance float64
	err = json.Unmarshal([]byte(val), &balance)
	return balance, err
}

func (ws *WalletService) Deposit(userID uuid.UUID, userTier string, amount float64) (*models.Transaction, error) {
	fmt.Println("userid: ", userID)
	if amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}

	balance, _ := ws.GetBalance(userID)
	newBalance := balance + amount

	err := ws.setBalance(userID, newBalance)
	if err != nil {
		return nil, err
	}

	txn := ws.createTransaction(userID, userTier, amount, models.DepositTransaction)
	ws.dbWriter <- txn

	return &txn, nil
}

func (ws *WalletService) Withdraw(userID uuid.UUID, userTier string, amount float64) (*models.Transaction, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}

	balance, _ := ws.GetBalance(userID)
	if balance < amount {
		return nil, errors.New("insufficient balance")
	}

	newBalance := balance - amount

	err := ws.setBalance(userID, newBalance)
	if err != nil {
		return nil, err
	}

	txn := ws.createTransaction(userID, userTier, amount, models.WithdrawTransaction)
	ws.dbWriter <- txn

	return &txn, nil
}

// Helper Functions
func (ws *WalletService) balanceKey(userID uuid.UUID) string {
	return fmt.Sprintf("wallet:balance:%s", userID)
}

func (ws *WalletService) setBalance(userID uuid.UUID, balance float64) error {
	// Store in Redis (optional but good for fast access)
	data, _ := json.Marshal(balance)
	if err := ws.redisClient.Set(ctx, ws.balanceKey(userID), data, 0).Err(); err != nil {
		return err
	}

	// Store in DB (upsert WalletBalance)
	walletBalance := models.WalletBalance{
		UserID:  userID,
		Balance: balance,
	}

	err := ws.db.
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"balance"}),
		}).
		Create(&walletBalance).Error

	return err
}

func (ws *WalletService) createTransaction(userID uuid.UUID, userTier string, amount float64, txnType models.TransactionType) models.Transaction {
	feeConfig := models.FeeConfig{TransactionType: "bill_payment", Tier: models.BasicTier, BasePercent: 3.0, Cap: 100.0, Floor: 2.0, PeakStart: time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 17, 0, 0, 0, time.Local), PeakEnd: time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 21, 0, 0, 0, time.Local), PeakSurcharge: 1.0}
	fee, breakdown := calculateFee(amount, userTier, feeConfig, time.Now())

	breakdownJSON, _ := json.Marshal(breakdown)

	return models.Transaction{
		UserID: userID,
		Amount: amount,
		Type:   txnType,
		// CreatedAt: time.Now(),
		Status:       "pending",
		Fee:          fee,
		NetAmount:    amount + fee,
		FeeBreakdown: breakdownJSON,
	}
}
func (ws *WalletService) GetTransactions(userID uuid.UUID, txnType string, status string, limit int) ([]models.Transaction, error) {
	var transactions []models.Transaction

	query := ws.db.Where("user_id = ?", userID)

	if txnType != "" {
		query = query.Where("type = ?", txnType)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if limit == 0 {
		limit = 50
	}

	err := query.Limit(limit).Order("created_at desc").Preload("User").Find(&transactions).Error
	return transactions, err
}

func calculateFee(amount float64, userTier string, config models.FeeConfig, now time.Time) (fee float64, breakdown map[string]interface{}) {
	// Base fee
	var basePercent float64
	if userTier == "Premium" {
		basePercent = 1

	} else {
		basePercent = 3
	}
	fee = amount * float64(basePercent) / 100

	// Time-based surcharge
	if now.After(config.PeakStart) && now.Before(config.PeakEnd) {
		surcharge := amount * config.PeakSurcharge / 100
		fee += surcharge
		breakdown = map[string]interface{}{
			"base_fee":       fee - surcharge,
			"peak_surcharge": surcharge,
		}
	} else {
		breakdown = map[string]interface{}{
			"base_fee":       fee,
			"peak_surcharge": 0,
		}
	}

	// Enforce floor and cap
	if fee < config.Floor {
		fee = config.Floor
	}
	if fee > config.Cap {
		fee = config.Cap
	}

	breakdown["total_fee"] = fee
	return fee, breakdown
}
