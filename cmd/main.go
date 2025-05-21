package main

import (
	"log"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/nazrawigedion123/wallet-backend/auth/handlers"
	auth_middleware "github.com/nazrawigedion123/wallet-backend/auth/middleware"
	"github.com/nazrawigedion123/wallet-backend/auth/services"
	db "github.com/nazrawigedion123/wallet-backend/utils"
)

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

	// Public routes
	e.POST("/register", authHandler.Register)
	e.POST("/login", authHandler.Login)

	e.Validator = &db.CustomValidator{Validator: validator.New()}

	// Protected routes
	authGroup := e.Group("")
	authGroup.Use(auth_middleware.AuthMiddleware(sessionSvc))
	authGroup.GET("/profile", authHandler.Profile)
	authGroup.POST("/tiers/upgrade", authHandler.TierUpgrade)
	authGroup.POST("/logout", authHandler.Logout)

	return e
}
