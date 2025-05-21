package middleware

import (
	"net/http"
	"strings"

	"github.com/nazrawigedion123/wallet-backend/auth/services"

	"github.com/labstack/echo/v4"
)

type AuthContext struct {
	echo.Context
	UserID       uint
	UserTier     string
	SessionToken string
}

func AuthMiddleware(sessionSvc *services.SessionService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "authorization header required"})
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "bearer token required"})
			}

			metadata, err := sessionSvc.ValidateSession(tokenString)

			if err != nil {

				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid session"})
			}

			// Store in context
			c.Set("userID", metadata.UserID)
			c.Set("userTier", metadata.Tier)
			c.Set("sessionToken", tokenString)

			return next(c)
		}
	}
}

// Helper to get auth context

func GetAuthContext(c echo.Context) *AuthContext {
	userIDValue := c.Get("userID")
	userTierValue := c.Get("userTier")
	sessionTokenValue := c.Get("sessionToken")

	userID, ok1 := userIDValue.(uint)
	userTier, ok2 := userTierValue.(string)
	sessionToken, ok3 := sessionTokenValue.(string)

	if !ok1 || !ok2 || !ok3 {
		return nil
	}

	return &AuthContext{
		Context:      c,
		UserID:       userID,
		UserTier:     userTier,
		SessionToken: sessionToken,
	}
}
