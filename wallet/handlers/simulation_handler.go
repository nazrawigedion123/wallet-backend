package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nazrawigedion123/wallet-backend/wallet/models"
)

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
