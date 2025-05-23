package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
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

// func ValidateHMACMiddleware() echo.MiddlewareFunc {
// 	secret := os.Getenv("WEBHOOK_SECRET")

// 	return func(next echo.HandlerFunc) echo.HandlerFunc {
// 		return func(c echo.Context) error {
// 			body, err := io.ReadAll(c.Request().Body)
// 			if err != nil {
// 				return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to read request body"})
// 			}

// 			c.Request().Body = io.NopCloser(bytes.NewBuffer(body)) // Reset body for handler

// 			// Compute HMAC
// 			mac := hmac.New(sha256.New, []byte(secret))

// 			mac.Write(body)
// 			fmt.Println(mac)
// 			fmt.Println("RAW BODY:\n", string(body))
// 			expectedSig := hex.EncodeToString(mac.Sum(nil))

// 			receivedSig := c.Request().Header.Get("X-Signature")
// 			fmt.Println(secret)
// 			fmt.Println("expected:", expectedSig)
// 			fmt.Println("received:", receivedSig)
// 			if receivedSig == "" || !hmac.Equal([]byte(receivedSig), []byte(expectedSig)) {
// 				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid or missing HMAC signature"})
// 			}

//				return next(c)
//			}
//		}
//	}
func ValidateHMACMiddleware() echo.MiddlewareFunc {
	secret := os.Getenv("WEBHOOK_SECRET")

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Read the request body
			body, err := io.ReadAll(c.Request().Body)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to read request body"})
			}

			// Reset the body for subsequent handlers
			c.Request().Body = io.NopCloser(bytes.NewBuffer(body))

			// Parse just the event_id from the JSON
			var payload struct {
				EventID string `json:"event_id"`
			}
			if err := json.Unmarshal(body, &payload); err != nil {
				return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid JSON"})
			}

			if payload.EventID == "" {
				return c.JSON(http.StatusBadRequest, echo.Map{"error": "event_id is required"})
			}

			// Compute HMAC using only the event_id value
			mac := hmac.New(sha256.New, []byte(secret))
			mac.Write([]byte(payload.EventID)) // Only use the event_id string
			expectedSig := hex.EncodeToString(mac.Sum(nil))

			receivedSig := c.Request().Header.Get("X-Signature")
			if receivedSig == "" || !hmac.Equal([]byte(receivedSig), []byte(expectedSig)) {
				return c.JSON(http.StatusUnauthorized, echo.Map{
					"error":    "invalid HMAC signature",
					"expected": expectedSig,
					"received": receivedSig,
					"message":  "HMAC is computed using only the event_id value",
				})
			}

			return next(c)
		}
	}
}
