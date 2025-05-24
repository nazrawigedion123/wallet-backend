package handlers

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nazrawigedion123/wallet-backend/webhook/models"
	"github.com/nazrawigedion123/wallet-backend/webhook/services"
)

type WebhookHandler struct {
	WebhookService *services.WebhookService
}

func NewWebhookHandler(webhookService *services.WebhookService) *WebhookHandler {
	return &WebhookHandler{
		WebhookService: webhookService,
	}
}

// HandleWebhook godoc
// @Summary Process incoming webhook
// @Description Handles and processes incoming webhook payloads
// @Tags Webhooks
// @Accept json
// @Produce json
// @Param payload body models.IncomingWebhook true "Webhook payload"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]interface{} "Returns error details"
// @Failure 409 {object} map[string]string "Duplicate webhook event"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /webhook [post]
func (h *WebhookHandler) HandleWebhook(c echo.Context) error {

	var payload models.IncomingWebhook
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid payload", "details": err.Error()})
	}
	if err := c.Validate(payload); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "validation failed", "details": err.Error()})
	}

	if err := h.WebhookService.ProcessWebhook(c.Request().Context(), payload); err != nil {
		if errors.Is(err, errors.New("duplicate webhook event")) {
			return c.JSON(409, map[string]string{"error": err.Error()})
		}
		return c.JSON(500, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"status": "received"})
}
