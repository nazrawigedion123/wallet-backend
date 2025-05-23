package routes

import (
	"github.com/labstack/echo/v4"

	"github.com/nazrawigedion123/wallet-backend/webhook/handlers"
	"github.com/nazrawigedion123/wallet-backend/webhook/middleware"
)

func RegisterWebhookRoutes(e *echo.Group, webhookHandler *handlers.WebhookHandler) {
	e.POST("/webhook/notify", webhookHandler.HandleWebhook, middleware.ValidateHMACMiddleware())
}
