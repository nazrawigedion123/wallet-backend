package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/nazrawigedion123/wallet-backend/auth/middleware"
	"github.com/nazrawigedion123/wallet-backend/auth/services"
	"github.com/nazrawigedion123/wallet-backend/wallet/handlers"
)

func RegisterWalletRoutes(e *echo.Group, walletHandler *handlers.WalletHandler, sessionSvc *services.SessionService) {
	walletGroup := e.Group("")
	walletGroup.Use(middleware.AuthMiddleware(sessionSvc))
	walletGroup.GET("/wallet/balance", walletHandler.GetBalance)
	walletGroup.POST("/wallet/deposit", walletHandler.Deposit)
	walletGroup.POST("/wallet/withdraw", walletHandler.Withdraw)
	walletGroup.GET("/wallet/transactions", walletHandler.GetTransactionHistory)
}

func RegisterSimulationRoutes(e *echo.Group, walletHandler *handlers.WalletHandler){
	simGroup := e.Group("")

	simGroup.POST("/simulate/users", walletHandler.SimulateUsers)
}
