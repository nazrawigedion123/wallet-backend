package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)
func GenerateHMACSignature(secret string, body []byte) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	return hex.EncodeToString(mac.Sum(nil))
}

func ValidateHMACMiddleware() echo.MiddlewareFunc {
	secret := os.Getenv("WEBHOOK_SECRET")

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			body, err := io.ReadAll(c.Request().Body)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to read request body"})
			}
			c.Request().Body = io.NopCloser(bytes.NewBuffer(body)) // Reset body for handler

			// Compute HMAC
			mac := hmac.New(sha256.New, []byte(secret))
			mac.Write(body)
			expectedSig := hex.EncodeToString(mac.Sum(nil))

			receivedSig := c.Request().Header.Get("X-Signature")
			if receivedSig == "" || !hmac.Equal([]byte(receivedSig), []byte(expectedSig)) {
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid or missing HMAC signature"})
			}

			return next(c)
		}
	}
}
