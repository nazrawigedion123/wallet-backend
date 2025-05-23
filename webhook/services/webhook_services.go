package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	transactionModels "github.com/nazrawigedion123/wallet-backend/wallet/models"
	"github.com/nazrawigedion123/wallet-backend/webhook/models"
)

type WebhookService struct {
	Redis *redis.Client
	DB    *gorm.DB
}

func NewWebhookService(redisClient *redis.Client, db *gorm.DB) *WebhookService {
	return &WebhookService{
		Redis: redisClient,
		DB:    db,
	}
}

func (s *WebhookService) ProcessWebhook(ctx context.Context, payload models.IncomingWebhook) error {
	idempotencyKey := fmt.Sprintf("webhook:event:%s", payload.EventID)

	// Idempotency check
	exists, err := s.Redis.Get(ctx, idempotencyKey).Result()
	if err == nil && exists == "1" {
		return errors.New("duplicate webhook event")
	}

	// Persist webhook event
	if err := s.saveWebhookEvent(ctx, payload); err != nil {
		return err
	}

	// Set idempotency key and payload in Redis
	if err := s.cacheWebhookEvent(ctx, idempotencyKey, payload); err != nil {
		return err
	}

	// Dispatch event
	switch payload.Type {
	case "wallet_credit":
		return s.handleWalletCredit(ctx, payload)
	case "wallet_debit":
		return s.handleWalletDebit(ctx, payload)
	case "bill_payment":
		return s.handleBillPayment(ctx, payload)
	default:
		return fmt.Errorf("unknown event type: %s", payload.Type)
	}
}

func (s *WebhookService) saveWebhookEvent(ctx context.Context, payload models.IncomingWebhook) error {
	event := models.WebhookEvent{
		EventID: payload.EventID,
		Type:    payload.Type,
		UserID:  payload.UserID,
		Amount:  payload.Amount,
	}
	if err := s.DB.WithContext(ctx).Create(&event).Error; err != nil {
		return fmt.Errorf("failed to save webhook event: %v", err)
	}
	return nil
}

func (s *WebhookService) cacheWebhookEvent(ctx context.Context, key string, payload models.IncomingWebhook) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %v", err)
	}

	pipe := s.Redis.TxPipeline()
	pipe.Set(ctx, key, "1", 24*time.Hour)
	pipe.Set(ctx, fmt.Sprintf("webhook:event:data:%s", payload.EventID), data, 24*time.Hour)

	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("failed to cache webhook event: %v", err)
	}
	return nil
}

func (s *WebhookService) handleWalletCredit(ctx context.Context, payload models.IncomingWebhook) error {
	tx := s.DB.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin DB transaction: %v", tx.Error)
	}

	// 1. Try to update the matching transaction first
	updateTx := &gorm.DB{} // declare updateTx to avoid undefined error

	if payload.Status == string(transactionModels.StatusSuccess) || payload.Status == string(transactionModels.StatusFailed) {
		updateTx = tx.Exec(`
		UPDATE transactions
		SET status = ?
		WHERE ctid IN (
			SELECT ctid FROM transactions
			WHERE user_id = ? AND amount = ? AND type = ? AND status = ?
			ORDER BY created_at ASC
			LIMIT 1
		)
	`, payload.Status, payload.UserID, payload.Amount, "deposit", transactionModels.StatusPending)
	} else {
		return fmt.Errorf("invalid type statys ")

	}

	if updateTx.Error != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update transaction status: %v", updateTx.Error)
	}
	if updateTx.RowsAffected == 0 {
		tx.Rollback()
		return fmt.Errorf("no matching pending transaction found to update")
	}

	// 2. Update wallet balance
	res := tx.Exec(`
		UPDATE wallet_balances
		SET balance = balance + ?
		WHERE user_id = ?`, payload.Amount, payload.UserID)

	if res.Error != nil {
		tx.Rollback()
		return fmt.Errorf("credit failed: %v", res.Error)
	}
	if res.RowsAffected == 0 {
		tx.Rollback()
		return fmt.Errorf("user wallet not found")
	}

	// 3. Commit
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	// 4. Update Redis balance
	go func() {
		pasrsedUUID, err := uuid.Parse(payload.UserID)
		if err != nil {
			fmt.Println("error parsing uid", err.Error())
		}
		_ = s.updateRedisBalance(pasrsedUUID)
	}()

	// 5. Publish Redis event
	s.Redis.Publish(ctx, "wallet:credit", fmt.Sprintf("user:%s:amount:%f", payload.UserID, payload.Amount))
	return nil
}

func (s *WebhookService) handleWalletDebit(ctx context.Context, payload models.IncomingWebhook) error {
	tx := s.DB.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin DB transaction: %v", tx.Error)
	}

	// 1. Check balance
	var balance float64
	err := tx.Raw(`
		SELECT balance FROM wallet_balances
		WHERE user_id = ?`, payload.UserID).Scan(&balance).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("balance check failed: %v", err)
	}
	if balance < payload.Amount {
		tx.Rollback()
		return fmt.Errorf("insufficient balance")
	}

	// 2. Try to update the transaction first
	updateTx := &gorm.DB{}
	// 1. Update wallet balance
	
	if payload.Status == string(transactionModels.StatusSuccess) || payload.Status == string(transactionModels.StatusFailed) {
		updateTx = tx.Exec(`
		UPDATE transactions
		SET status = ?
		WHERE ctid IN (
			SELECT ctid FROM transactions
			WHERE user_id = ? AND amount = ? AND type = ? AND status = ?
			ORDER BY created_at ASC
			LIMIT 1
		)
	`, payload.Status, payload.UserID, payload.Amount, "withdraw", transactionModels.StatusPending)

	} else {
		return fmt.Errorf("invalid payload type")
	}

	if updateTx.Error != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update transaction status: %v", updateTx.Error)
	}
	if updateTx.RowsAffected == 0 {
		tx.Rollback()
		return fmt.Errorf("no matching pending withdrawal transaction found to update")
	}

	// 3. Deduct from wallet
	res := tx.Exec(`
		UPDATE wallet_balances
		SET balance = balance - ?
		WHERE user_id = ?`, payload.Amount, payload.UserID)

	if res.Error != nil {
		tx.Rollback()
		return fmt.Errorf("debit failed: %v", res.Error)
	}
	if res.RowsAffected == 0 {
		tx.Rollback()
		return fmt.Errorf("user wallet not found")
	}

	// 4. Commit
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	// 5. Update Redis balance
	go func() {
		pasrsedUUID, err := uuid.Parse(payload.UserID)
		if err != nil {
			fmt.Println("error parsing uid", err.Error())
		}
		_ = s.updateRedisBalance(pasrsedUUID)
	}()

	// 6. Publish Redis event
	s.Redis.Publish(ctx, "wallet:debit", fmt.Sprintf("user:%s:amount:%f", payload.UserID, payload.Amount))
	return nil
}

func (s *WebhookService) updateRedisBalance(userID uuid.UUID) error {
	var balance float64
	err := s.DB.
		Raw("SELECT balance FROM wallet_balances WHERE user_id = ?", userID).
		Scan(&balance).Error
	if err != nil {
		return err
	}

	data, _ := json.Marshal(balance)
	return s.Redis.Set(context.Background(), fmt.Sprintf("wallet:balance:%s", userID), data, 0).Err()
}

func (s *WebhookService) handleBillPayment(ctx context.Context, payload models.IncomingWebhook) error {
	if err := s.handleWalletDebit(ctx, payload); err != nil {
		return err
	}
	// Optional: Save bill payment record or emit event
	return nil
}
