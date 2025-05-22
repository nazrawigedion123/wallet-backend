package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/nazrawigedion123/wallet-backend/auth/handlers"
	"github.com/nazrawigedion123/wallet-backend/auth/middleware"
	"github.com/nazrawigedion123/wallet-backend/auth/services"
)

func RegisterAuthRoutes(e *echo.Group, authHandler *handlers.AuthHandler, sessionSvc *services.SessionService) {
	// Public
	e.POST("/register", authHandler.Register)
	e.POST("/login", authHandler.Login)

	// Protected
	authGroup := e.Group("")
	authGroup.Use(middleware.AuthMiddleware(sessionSvc))
	authGroup.GET("/profile", authHandler.Profile)
	authGroup.POST("/tiers/upgrade", authHandler.TierUpgrade)
	authGroup.POST("/logout", authHandler.Logout)
}
