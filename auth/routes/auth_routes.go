package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/nazrawigedion123/wallet-backend/auth/handlers"
	"github.com/nazrawigedion123/wallet-backend/auth/middleware"
	"github.com/nazrawigedion123/wallet-backend/auth/services"
)

func RegisterAuthRoutes(e *echo.Group, authHandler *handlers.AuthHandler, sessionSvc *services.SessionService) {
	// Public

	// Protected
	authGroup := e.Group("")
	authGroup.POST("/register", authHandler.Register)
	authGroup.POST("/login", authHandler.Login)
	authGroup.Use(middleware.AuthMiddleware(sessionSvc))
	authGroup.GET("/profile", authHandler.Profile)
	authGroup.POST("/tiers/upgrade", authHandler.TierUpgrade)
	authGroup.POST("/logout", authHandler.Logout)
}
