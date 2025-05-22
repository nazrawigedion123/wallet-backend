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

func (h *WalletHandler) Deposit(c echo.Context) error {
	userID := c.Get("userID").(uuid.UUID)

	var req TransactionRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	txn, err := h.WalletService.Deposit(userID, req.Amount)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, txn)
}

func (h *WalletHandler) Withdraw(c echo.Context) error {

	userID := c.Get("userID").(uuid.UUID)

	var req TransactionRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	txn, err := h.WalletService.Withdraw(userID, req.Amount)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, txn)
}
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
