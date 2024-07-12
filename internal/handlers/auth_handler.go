package handlers

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	dto "github.com/mohamed2394/sahla/internal/dtos"
	service "github.com/mohamed2394/sahla/internal/services"

)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Login(c echo.Context) error {
	var loginRequest dto.LoginRequest
	if err := c.Bind(&loginRequest); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid request format"})
	}

	// Validate the request
	if err := c.Validate(loginRequest); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
	}

	// Attempt to authenticate and get a token
	accesToken, _, err := h.authService.Login(loginRequest.Email, loginRequest.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "Invalid credentials"})
	}

	// Return the token
	return c.JSON(http.StatusOK, dto.LoginResponse{Token: accesToken})
}

func (h *AuthHandler) Logout(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")
	if token == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "No token provided"})
	}

	// Remove "Bearer " prefix if present
	token = strings.TrimPrefix(token, "Bearer ")

	err := h.authService.Logout(token)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to logout"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Successfully logged out"})
}

func (h *AuthHandler) RefreshToken(c echo.Context) error {
	// Get the refresh token from the request
	refreshToken := c.FormValue("refresh_token")
	if refreshToken == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Refresh token is required"})
	}

	// Call the service to handle token refresh
	newAccessToken, newRefreshToken, err := h.authService.RefreshToken(refreshToken)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	// Return the new tokens
	return c.JSON(http.StatusOK, map[string]string{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	})
}
