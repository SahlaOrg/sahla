package routes

import (
	"github.com/labstack/echo/v4"
	handler "github.com/mohamed2394/sahla/internal/handlers"
)

func RegisterCreditPaymentRoutes(e *echo.Echo, creditPaymentHandler *handler.CreditPaymentHandler) {
	e.POST("/credit-applications", creditPaymentHandler.CreateCreditApplication)
	e.PUT("/credit-applications/:id/approve", creditPaymentHandler.ApproveCreditApplication)
	e.POST("/payments", creditPaymentHandler.CreatePayment)
	e.GET("/payments/:id", creditPaymentHandler.GetPaymentDetails)
	e.POST("/installments/:id/process", creditPaymentHandler.ProcessInstallment)
	e.POST("/webhooks/payments/:id", creditPaymentHandler.HandlePaymentWebhook)
	e.POST("/webhooks/installments/:id", creditPaymentHandler.HandleInstallmentWebhook)
}