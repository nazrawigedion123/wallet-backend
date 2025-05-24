package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/nazrawigedion123/wallet-backend/wallet/services"
)

type WalletHandler struct {
	WalletService *services.WalletService
}

type TransactionRequest struct {
	Amount float64 `json:"amount" validate:"required,gt=0"`
}

// GetBalance godoc
// @Summary Get user wallet balance
// @Description Returns the wallet balance for the authenticated user
// @Tags Wallet
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /wallet/balance [get]
func (h *WalletHandler) GetBalance(c echo.Context) error {
	fmt.Println("Get balance handler")
	userID := c.Get("userID").(uuid.UUID)

	balance, err := h.WalletService.GetBalance(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not get balance")
	}

	return c.JSON(http.StatusOK, echo.Map{
		"user_id": userID,
		"balance": balance,
	})
}
// Deposit godoc
// @Summary Deposit to wallet 
// @Description it deposits a to transaction and adds money to balance
// @Tags Wallet
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /wallet/deposit [post]
func (h *WalletHandler) Deposit(c echo.Context) error {
	userID := c.Get("userID").(uuid.UUID)
	userTierInterface := c.Get("userTier")
	userTier, ok := userTierInterface.(string)
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid user tier")
	}

	var req TransactionRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	txn, err := h.WalletService.Deposit(userID, userTier, req.Amount)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, txn)
}

// Withdraw godoc
// @Summary Withdraw from wallet 
// @Description Withdraws money from the wallet balance
// @Tags Wallet
// @Security BearerAuth
// @Produce json
// @Param request body TransactionRequest true "Withdrawal request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /wallet/withdraw [post]
func (h *WalletHandler) Withdraw(c echo.Context) error {

	userID := c.Get("userID").(uuid.UUID)
	userTierInterface := c.Get("userTier")
	userTier, ok := userTierInterface.(string)
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid user tier")
	}


	var req TransactionRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	txn, err := h.WalletService.Withdraw(userID,userTier, req.Amount)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, txn)
}

// GetTransactionHistory godoc
// @Summary Get transaction history
// @Description Retrieves the transaction history for the user's wallet
// @Tags Wallet
// @Security BearerAuth
// @Produce json
// @Param type query string false "Filter by transaction type"
// @Param status query string false "Filter by transaction status"
// @Param limit query integer false "Limit number of transactions (default 50)"
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /wallet/transactions [get]
func (h *WalletHandler) GetTransactionHistory(c echo.Context) error {
	userID := c.Get("userID").(uuid.UUID)

	txnType := c.QueryParam("type")
	status := c.QueryParam("status")

	limit := 50
	if l := c.QueryParam("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}

	transactions, err := h.WalletService.GetTransactions(userID, txnType, status, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "could not fetch transactions"})
	}

	return c.JSON(http.StatusOK, transactions)
}
