// internal/handler/payment_handler.go
package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/mohamed2394/sahla/modules/credit/dto"
	"github.com/mohamed2394/sahla/modules/credit/service"
)

type PaymentHandler struct {
	paymentService service.PaymentService
}

func NewPaymentHandler(paymentService service.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
	}
}

func (h *PaymentHandler) MakePayment(c echo.Context) error {
	var req dto.CreatePaymentRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	payment, err := h.paymentService.MakePayment(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, payment)
}

func (h *PaymentHandler) GetPaymentsByLoanID(c echo.Context) error {
	loanID, err := strconv.ParseUint(c.Param("loanID"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid loan ID"})
	}

	payments, err := h.paymentService.GetPaymentsByLoanID(c.Request().Context(), uint(loanID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, payments)
}
