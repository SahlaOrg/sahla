package routes

import (
	"github.com/labstack/echo/v4"
	handler "github.com/mohamed2394/sahla/modules/credit/handler"
)

func RegisterCreditRoutes(e *echo.Echo, creditHandler *handler.CreditHandler) {
	e.POST("/credit/score", creditHandler.SaveCreditScore)
	e.GET("/credit/score/:user_id", creditHandler.GetCreditScoreByUserID)
	e.POST("/credit/features", creditHandler.SaveCreditFeatures)
	e.GET("/credit/features/:user_id", creditHandler.GetCreditFeaturesByUserID)
}
