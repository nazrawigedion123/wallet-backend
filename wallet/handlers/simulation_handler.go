package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nazrawigedion123/wallet-backend/wallet/models"
	"github.com/nazrawigedion123/wallet-backend/wallet/utils"
)
// SimulateUsers godoc
// @Summary Simulate user creation
// @Description Starts a background process to generate simulated users
// @Tags Simulation
// @Security BearerAuth
// @Produce json
// @Param request body models.SimulationOptions true "Simulation options"
// @Success 202 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /simulate/users [post]
func (h *WalletHandler) SimulateUsers(c echo.Context) error {
	var opts models.SimulationOptions
	if err := c.Bind(&opts); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}
	// Validate count
	if opts.Count <= 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Count must be greater than 0"})
	}

	go h.WalletService.GenerateUsers(opts) // run in background

	return c.JSON(http.StatusAccepted, echo.Map{"message": "User simulation started"})
}

// SimulateTransactions godoc
// @Summary Simulate transactions
// @Description Starts a background process to generate simulated transactions
// @Tags Simulation
// @Security BearerAuth
// @Produce json
// @Param request body models.SimulationOptions true "Simulation options"
// @Success 202 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /simulate/transactions [post]
func (h *WalletHandler) SimulateTransactions(c echo.Context) error {
	var opts models.SimulationOptions
	if err := c.Bind(&opts); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}

	if opts.Count <= 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Count must be > 0"})
	}

	go h.WalletService.GenerateTransactions(opts)

	return c.JSON(http.StatusAccepted, echo.Map{
		"message": "Transaction simulation started",
		"count":   opts.Count,
	})
}
// GetSimulationStatus godoc
// @Summary Get simulation status
// @Description Returns the current status of background simulations
// @Tags Simulation
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /simulate/status [get]
func (h *WalletHandler) GetSimulationStatus(c echo.Context) error {
	statuses := utils.GetStatuses()
	return c.JSON(http.StatusOK, statuses)
}


