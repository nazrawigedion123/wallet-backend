// @title Wallet Backend API
// @version 1.0
// @description This is a wallet backend server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
package main

import (
	"log"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/nazrawigedion123/wallet-backend/auth/handlers"
	"github.com/nazrawigedion123/wallet-backend/auth/services"
	db "github.com/nazrawigedion123/wallet-backend/utils"

	_ "github.com/nazrawigedion123/wallet-backend/docs"
	echoSwagger "github.com/swaggo/echo-swagger"

	authRoutes "github.com/nazrawigedion123/wallet-backend/auth/routes"
	walletHandler "github.com/nazrawigedion123/wallet-backend/wallet/handlers"
	walletRoutes "github.com/nazrawigedion123/wallet-backend/wallet/routes"
	walletService "github.com/nazrawigedion123/wallet-backend/wallet/services"
)

var redisClient *redis.Client

// @Summary Register a new user
// @Description Register a new user with the system
// @ID register
// @Accept  json
// @Produce  json
func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found or failed to load")
	}
	if err := initDatabase(); err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	defer db.CloseConnections()

	if err := initRedis(); err != nil {
		log.Fatalf("❌ Failed to connect to Redis: %v", err)
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("❌ JWT_SECRET environment variable is not set")
	}

	//auth
	sessionSvc, authSvc := initServices(jwtSecret)
	authHandler := handlers.NewAuthHandler(authSvc, sessionSvc)

	e := setupServer(authHandler, sessionSvc)

	log.Println("🚀 Server started on :8080")
	e.Logger.Fatal(e.Start(":8080"))
}

func initDatabase() error {
	return db.InitDB()
}

func initRedis() error {
	config := db.RedisConfig{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	}
	redisClient = redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       config.DB,
	})

	return db.InitRedis()
}

func initServices(jwtSecret string) (*services.SessionService, *services.AuthService) {
	sessionSvc := services.NewSessionService(db.RedisClient, jwtSecret, 24*time.Hour)
	authSvc := services.NewAuthService(db.DB, sessionSvc)
	return sessionSvc, authSvc
}

func setupServer(authHandler *handlers.AuthHandler, sessionSvc *services.SessionService) *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Add Swagger route
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	apiGroup := e.Group("/api")
	e.Validator = &db.CustomValidator{Validator: validator.New()}

	authRoutes.RegisterAuthRoutes(apiGroup, authHandler, sessionSvc)

	ws := walletService.NewWalletService(db.DB, redisClient)
	walletHandlerInstance := &walletHandler.WalletHandler{
		WalletService: ws,
	}

	walletRoutes.RegisterWalletRoutes(apiGroup, walletHandlerInstance, sessionSvc)

	return e
}
