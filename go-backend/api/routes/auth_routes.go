package routes

import (
	"github.com/labstack/echo/v4"
	handler "github.com/mohamed2394/sahla/modules/auth/handler"
	"github.com/mohamed2394/sahla/modules/auth/service"
	"github.com/mohamed2394/sahla/modules/user/repository"
)

func SetupAuthRoutes(e *echo.Echo, userRepo repository.UserRepository, jwtSecret, refreshSecret string) {
	// Initialize services
	authService := service.NewAuthService(userRepo, jwtSecret, refreshSecret)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)

	// Auth routes
	e.POST("/login", authHandler.Login)
	// e.POST("/refresh", authHandler.RefreshToken)
	e.POST("/logout", authHandler.Logout)
}
