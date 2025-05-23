package utils

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/go-playground/validator/v10"

	user_models "github.com/nazrawigedion123/wallet-backend/auth/models"
	wallet_models "github.com/nazrawigedion123/wallet-backend/wallet/models"
	webhook_models "github.com/nazrawigedion123/wallet-backend/webhook/models"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}

var (
	DB          *gorm.DB
	RedisClient *redis.Client
)

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

func InitDB() error {

	config := DBConfig{
		Host:     os.Getenv("HOST"),
		Port:     os.Getenv("PORT"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("PASSWORD"),
		DBName:   os.Getenv("DBNAME"),
	}
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.DBName)
	fmt.Println("dsn: ", dsn)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	// Auto-migrate the Transaction model
	err = DB.AutoMigrate(&user_models.User{},
		&wallet_models.Transaction{},
		&wallet_models.WalletBalance{},
		&webhook_models.WebhookEvent{},
	)
	if err != nil {
		return err
	}

	log.Println("Auto migration complete")

	log.Println("Connected to PostgreSQL database")
	return nil
}

func InitRedis() error {
	// Initialize Redis
	config := RedisConfig{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	}
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       config.DB,
	})

	ctx := context.Background()
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		return err
	}

	log.Println("Connected to Redis")
	return nil
}

func CloseConnections() {
	if sqlDB, err := DB.DB(); err == nil {
		sqlDB.Close()
	}
	if RedisClient != nil {
		RedisClient.Close()
	}
}
