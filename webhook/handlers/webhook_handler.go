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

// func (h *WebhookHandler) HandleWebhook(c echo.Context) error {
// 	var payload models.IncomingWebhook

// 	if err := c.Bind(&payload); err != nil {
// 		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid payload", "details": err.Error()})
// 	}

// 	if err := c.Validate(payload); err != nil {
// 		return c.JSON(http.StatusBadRequest, echo.Map{"error": "validation failed", "details": err.Error()})
// 	}

// 	if err := h.WebhookService.ProcessWebhook(c.Request().Context(), payload); err != nil {
// 		// Check error type for appropriate status code
// 		if errors.Is(err, errors.New("duplicate webhook event")) {
// 			return c.JSON(409, map[string]string{"error": err.Error()})
// 		}
// 		return c.JSON(500, map[string]string{"error": err.Error()})
// 	}

//		return c.JSON(http.StatusOK, echo.Map{"status": "received"})
//	}
func (h *WebhookHandler) HandleWebhook(c echo.Context) error {

	// // Read raw body for HMAC verification
	// body, err := io.ReadAll(c.Request().Body)
	// if err != nil {
	// 	return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to read request body"})
	// }
	// // Restore body for Bind
	// c.Request().Body = io.NopCloser(bytes.NewBuffer(body))

	// // Compute expected signature
	// mac := hmac.New(sha256.New, []byte(secret))
	// mac.Write(body)
	// expectedSig := hex.EncodeToString(mac.Sum(nil))

	// // Compare with the signature sent in the header
	// receivedSig := c.Request().Header.Get("X-Signature")
	// fmt.Println("expected:", expectedSig)
	// fmt.Println("recived: ", receivedSig)
	// if receivedSig == "" || !hmac.Equal([]byte(receivedSig), []byte(expectedSig)) {
	// 	return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid HMAC signature"})

	// }

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
