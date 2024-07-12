package routes

import (
	"github.com/labstack/echo/v4"
	handler "github.com/mohamed2394/sahla/internal/handlers"

)

func RegisterUserRoutes(e *echo.Echo, userHandler *handler.UserHandler) {
	e.POST("/users", userHandler.CreateUser)
	e.GET("/users/:id", userHandler.GetUserByID)
	e.PUT("/users/:id", userHandler.UpdateUser)
	e.DELETE("/users/:id", userHandler.DeleteUser)
	e.GET("/users", userHandler.ListUsers)
	e.POST("/users/:id/upload-id-image", userHandler.UploadIDImage)

}
