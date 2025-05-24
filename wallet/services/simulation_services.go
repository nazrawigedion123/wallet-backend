package services

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/google/uuid"
	userModel "github.com/nazrawigedion123/wallet-backend/auth/models"
	"github.com/nazrawigedion123/wallet-backend/wallet/models"
	"gorm.io/gorm"
	"github.com/nazrawigedion123/wallet-backend/wallet/utils"
)

func (ws *WalletService) GenerateUsers(opts models.SimulationOptions) {
	rateLimiter := time.Tick(100 * time.Microsecond) // simple rate limit
	tiers := []string{"basic", "premium", "enterprise"}

	var records [][]string
	for i := 0; i < opts.Count; i++ {
		<-rateLimiter
		user := userModel.User{
			ID:    uuid.New(),
			Email: fmt.Sprintf("user%d@example.com", i),
			Tier:  tiers[rand.Intn(len(tiers))],
		}

		ws.db.Create(&user)

		if opts.OutputToCSV {
			records = append(records, []string{user.ID.String(), user.Email, user.Tier})
		}
		utils.IncrementProgress("users")
	}
	utils.EndSimulation("users")

	if opts.OutputToCSV {
		file, _ := os.Create("users_simulated.csv")
		writer := csv.NewWriter(file)
		writer.WriteAll(records)
		writer.Flush()
		file.Close()
	}

	log.Printf("✅ Finished simulating %d users", opts.Count)
}

func (ws *WalletService) GenerateTransactions(opts models.SimulationOptions) {
	rateLimiter := time.Tick(100 * time.Microsecond)

	var users []userModel.User
	ws.db.Find(&users)

	if len(users) == 0 {
		log.Println("❌ No users found to attach transactions")
		return
	}

	var records [][]string

	for i := 0; i < opts.Count; i++ {
		<-rateLimiter

		user := users[rand.Intn(len(users))]
		amount := rand.Float64()*1000 + 10 // between 10 and 1010
		txnType := models.TransactionType("payment")
		if len(opts.TransactionTypes) > 0 {
			txnType = models.TransactionType(opts.TransactionTypes[rand.Intn(len(opts.TransactionTypes))])
		}

		status := models.StatusPending
		if opts.IncludeFailed && rand.Intn(10) < 2 { // 20% chance to fail
			status = models.StatusFailed
		}

		fee, breakdown := calculateFee(amount, user.Tier, models.FeeConfig{}, time.Now())
		breakdownJSON, _ := json.Marshal(breakdown)

		txn := models.Transaction{
			UserID:       user.ID,
			Amount:       amount,
			Type:         txnType,
			Status:       status,
			Fee:          fee,
			NetAmount:    amount - fee,
			FeeBreakdown: breakdownJSON,
		}

		ws.db.Create(&txn)

		if status == models.StatusPending || status == models.StatusSuccess {
			ws.db.Model(&models.WalletBalance{}).
				Where("user_id = ?", user.ID).
				Update("balance", gorm.Expr("balance + ?", amount-fee))
		}

		if opts.OutputToCSV {
			records = append(records, []string{
				fmt.Sprintf("%d", txn.ID), user.Email, fmt.Sprintf("%.2f", amount),
				string(txnType), string(status),
				fmt.Sprintf("%.2f", fee), fmt.Sprintf("%.2f", amount-fee),
			})
		}
		utils.IncrementProgress("transactions")
	}
	utils.EndSimulation("transactions")

	if opts.OutputToCSV {
		file, _ := os.Create("transactions_simulated.csv")
		writer := csv.NewWriter(file)
		writer.WriteAll(records)
		writer.Flush()
		file.Close()
	}

	log.Printf("✅ Finished simulating %d transactions", opts.Count)
}
