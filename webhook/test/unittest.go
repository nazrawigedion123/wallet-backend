package test

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/nazrawigedion123/wallet-backend/webhook/handlers"
	"github.com/nazrawigedion123/wallet-backend/webhook/interfaces"
	"github.com/nazrawigedion123/wallet-backend/webhook/mocks"
	"github.com/nazrawigedion123/wallet-backend/webhook/models"
	"github.com/stretchr/testify/assert"
)

func GenerateHMACSignature(secret string, body []byte) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	return hex.EncodeToString(mac.Sum(nil))
}

func TestHandleWebhook_Success(t *testing.T) {
	mockService := &mocks.MockWebhookService{
		ProcessFunc: func(ctx context.Context, payload models.IncomingWebhook) error {
			return nil // Simulate success
		},
	}

	handler := handlers.NewWebhookHandler(mockService)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{
				"type": "wallet_credit",
				"event_id": "txn_acodnnnjnn",         
				"amount": 150,
				"currency": "ETB",
				"status": "success",
					"user_id": "95dfe0d1-0e5a-4da4-8dd8-f8afe9266f4a",
				"metadata": {
					"transaction_id": 5             
				}
				}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	// Act
	err := handler.HandleWebhook(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}
