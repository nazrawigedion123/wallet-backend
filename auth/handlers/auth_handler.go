package handlers

import (
	"net/http"

	"github.com/nazrawigedion123/wallet-backend/auth/models"

	"github.com/nazrawigedion123/wallet-backend/auth/middleware"
	"github.com/nazrawigedion123/wallet-backend/auth/services"
	"github.com/nazrawigedion123/wallet-backend/utils"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authSvc    *services.AuthService
	sessionSvc *services.SessionService
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=3"`
}
type TierUpgrade struct {
	Tier string `json:"tier" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func NewAuthHandler(authSvc *services.AuthService, sessionSvc *services.SessionService) *AuthHandler {
	return &AuthHandler{
		authSvc:    authSvc,
		sessionSvc: sessionSvc,
	}
}

func (h *AuthHandler) Register(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	user, err := h.authSvc.Register(req.Email, req.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "registration failed"})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"id":    user.ID,
		"email": user.Email,
		"tier":  user.Tier,
	})
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	ipAddress := c.RealIP()

	token, user, err := h.authSvc.Login(req.Email, req.Password, ipAddress)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
		"user": map[string]interface{}{
			"id":    user.ID,
			"email": user.Email,
			"tier":  user.Tier,
		},
	})
}

func (h *AuthHandler) Logout(c echo.Context) error {
	cc := middleware.GetAuthContext(c)
	if cc == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "not authenticated"})
	}

	err := h.sessionSvc.InvalidateSession(cc.SessionToken)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "logout failed"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "logged out successfully"})
}

func (h *AuthHandler) Profile(c echo.Context) error {

	userIDVal := c.Get("userID")

	if userIDVal == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "not authenticated"})
	}

	userID, ok := userIDVal.(uint)
	if !ok {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "invalid user id type"})
	}

	var user models.User
	if err := utils.DB.First(&user, userID).Error; err != nil {

		return c.JSON(http.StatusNotFound, map[string]string{"error": "user not found"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"id":    user.ID,
		"email": user.Email,
		"tier":  user.Tier,
	})
}

func (h *AuthHandler) TierUpgrade(c echo.Context) error {

	userIDVal := c.Get("userID")

	if userIDVal == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "not authenticated"})
	}

	userID, ok := userIDVal.(uint)
	if !ok {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "invalid user id type"})
	}
	var tier TierUpgrade
	if err := c.Bind(&tier); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	var user models.User
	if err := utils.DB.First(&user, userID).Error; err != nil {

		return c.JSON(http.StatusNotFound, map[string]string{"error": "user not found"})
	}
	if tier.Tier == "Premium" || tier.Tier == "Enterprise" || tier.Tier == "Basic" {
		user.Tier = tier.Tier

	} else {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid Tier"})

	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"id":    user.ID,
		"email": user.Email,
		"tier":  user.Tier,
	})
}
