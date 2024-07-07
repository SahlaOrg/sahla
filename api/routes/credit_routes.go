// internal/router/router.go
package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/mohamed2394/sahla/modules/credit/handler"
	"github.com/mohamed2394/sahla/modules/credit/service"
)

func SetupCreditRoutes(e *echo.Echo, loanService service.LoanService, paymentService service.PaymentService, repaymentPlanService service.RepaymentPlanService) {
	loanHandler := handler.NewLoanHandler(loanService)
	paymentHandler := handler.NewPaymentHandler(paymentService)
	repaymentPlanHandler := handler.NewRepaymentPlanHandler(repaymentPlanService)

	// Loan routes
	e.POST("/loans", loanHandler.CreateLoan)
	e.GET("/loans/:id", loanHandler.GetLoanByID)
	e.GET("/loans/universal/:universalID", loanHandler.GetLoanByUniversalID)
	e.GET("/loans/user/:userID", loanHandler.GetLoansByUserID)
	e.PUT("/loans/:id/status", loanHandler.UpdateLoanStatus)

	// Payment routes
	e.POST("/payments", paymentHandler.MakePayment)
	e.GET("/payments/loan/:loanID", paymentHandler.GetPaymentsByLoanID)

	// Repayment plan routes
	e.GET("/repayment-plans/loan/:loanID", repaymentPlanHandler.GetRepaymentPlanByLoanID)
}
