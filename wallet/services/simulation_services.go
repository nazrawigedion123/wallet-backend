package services

import (
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/nazrawigedion123/wallet-backend/wallet/models"
	userModel	"github.com/nazrawigedion123/wallet-backend/auth/models"
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
	}

	if opts.OutputToCSV {
		file, _ := os.Create("users_simulated.csv")
		writer := csv.NewWriter(file)
		writer.WriteAll(records)
		writer.Flush()
		file.Close()
	}

	log.Printf("✅ Finished simulating %d users", opts.Count)
}
