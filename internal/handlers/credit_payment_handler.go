package handlers

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	dto "github.com/mohamed2394/sahla/internal/dtos"
	domains "github.com/mohamed2394/sahla/internal/domains"
	services"github.com/mohamed2394/sahla/internal/services"
	utils "github.com/mohamed2394/sahla/internal/utils"
	validation"github.com/mohamed2394/sahla/internal/validation"
	"go.uber.org/zap"
)

// CreditPaymentHandler handles HTTP requests related to credit payments
type CreditPaymentHandler struct {
	service   services.CreditPaymentServiceInterface
	logger    *zap.Logger
	validator *validation.CustomValidator
}

// NewCreditPaymentHandler creates a new instance of CreditPaymentHandler
func NewCreditPaymentHandler(service services.CreditPaymentServiceInterface, logger *zap.Logger, validator *validation.CustomValidator) *CreditPaymentHandler {
	return &CreditPaymentHandler{
		service:   service,
		logger:    logger,
		validator: validator,
	}
}

// CreateCreditApplication handles the creation of a new credit application
func (h *CreditPaymentHandler) CreateCreditApplication(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	var req dto.CreditApplicationRequest
	if err := c.Bind(&req); err != nil {
		h.logger.Error("Failed to bind request body", zap.Error(err))
		return h.handleError(c, err, "invalid request body")
	}

	h.logger.Info("CreditApplicationRequest received", zap.Any("request", req))

	if err := h.validator.Validate(req); err != nil {
		h.logger.Error("Validation failed", zap.Error(err))
		return h.handleError(c, err, "validation failed")
	}

	if !isValidCurrency(req.Currency) {
		err := errors.New("unsupported currency")
		h.logger.Error("Validation failed", zap.Error(err))
		return h.handleError(c, err, "unsupported currency")
	}

	app := &domains.CreditApplication{
		UserID:   req.UserID,
		Amount:   req.Amount,
		Currency: req.Currency,
	}

	if err := h.service.CreateCreditApplication(ctx, app); err != nil {
		h.logger.Error("Failed to create credit application", zap.Error(err))
		return h.handleError(c, err, "failed to create credit application")
	}

	h.logger.Info("Credit application created successfully", zap.Any("application", app))
	return c.JSON(http.StatusCreated, h.createCreditApplicationResponse(app))
}

func isValidCurrency(currency string) bool {
	// Add supported currency validation logic here
	return currency == "DZD" // Example
}
// ApproveCreditApplication handles the approval of a credit application
func (h *CreditPaymentHandler) ApproveCreditApplication(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	id, err := h.parseID(c.Param("id"))
	if err != nil {
		h.logger.Error("Invalid credit application ID", zap.Error(err))
		return h.handleError(c, err, "invalid credit application ID")
	}

	h.logger.Info("Approving credit application", zap.Uint("applicationID", id))

	if err := h.service.ApproveCreditApplication(ctx, id); err != nil {
		h.logger.Error("Failed to approve credit application", zap.Error(err))
		return h.handleError(c, err, "failed to approve credit application")
	}

	h.logger.Info("Credit application approved successfully", zap.Uint("applicationID", id))
	return c.JSON(http.StatusOK, map[string]string{"message": "Credit application approved successfully"})
}
// CreatePayment handles the creation of a new payment
func (h *CreditPaymentHandler) CreatePayment(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	var req dto.PaymentRequest
	if err := c.Bind(&req); err != nil {
		h.logger.Error("Failed to bind request body", zap.Error(err))
		return h.handleError(c, err, "invalid request body")
	}

	h.logger.Info("PaymentRequest received", zap.Any("request", req))

	if err := h.validator.Validate(req); err != nil {
		h.logger.Error("Validation failed", zap.Error(err))
		return h.handleError(c, err, "validation failed")
	}

	if !isValidPaymentMethod(req.PaymentMethod) {
		err := errors.New("unsupported payment method")
		h.logger.Error("Validation failed", zap.Error(err))
		return h.handleError(c, err, "unsupported payment method")
	}

	payment := &domains.Payment{
		CreditApplicationID: req.CreditApplicationID,
		UserID:              req.UserID,
		Amount:              req.Amount,
		Currency:            req.Currency,
		PaymentMethod:       req.PaymentMethod,
	}

	if err := h.service.CreatePayment(ctx, payment); err != nil {
		h.logger.Error("Failed to create payment", zap.Error(err))
		return h.handleError(c, err, "failed to create payment")
	}

	h.logger.Info("Payment created successfully", zap.Any("payment", payment))
	return c.JSON(http.StatusCreated, h.createPaymentResponse(payment))
}

func isValidPaymentMethod(paymentMethod string) bool {
	// Add supported payment method validation logic here
	return paymentMethod == "card" // Example
}
// GetPaymentDetails retrieves the details of a specific payment
func (h *CreditPaymentHandler) GetPaymentDetails(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	id, err := h.parseID(c.Param("id"))
	if err != nil {
		h.logger.Error("Invalid payment ID", zap.Error(err))
		return h.handleError(c, err, "invalid payment ID")
	}

	h.logger.Info("Fetching payment details", zap.Uint("paymentID", id))

	payment, err := h.service.GetPaymentDetails(ctx, id)
	if err != nil {
		h.logger.Error("Failed to get payment details", zap.Error(err))
		return h.handleError(c, err, "failed to get payment details")
	}

	h.logger.Info("Payment details retrieved successfully", zap.Any("payment", payment))
	return c.JSON(http.StatusOK, h.createPaymentResponse(payment))
}
// ProcessInstallment handles the processing of an installment
func (h *CreditPaymentHandler) ProcessInstallment(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	id, err := h.parseID(c.Param("id"))
	if err != nil {
		h.logger.Error("Invalid installment ID", zap.Error(err))
		return h.handleError(c, err, "invalid installment ID")
	}

	h.logger.Info("Processing installment", zap.Uint("installmentID", id))

	if err := h.service.ProcessInstallment(ctx, id); err != nil {
		h.logger.Error("Failed to process installment", zap.Error(err))
		return h.handleError(c, err, "failed to process installment")
	}

	h.logger.Info("Installment processed successfully", zap.Uint("installmentID", id))
	return c.JSON(http.StatusOK, map[string]string{"message": "Installment processed successfully"})
}

// HandlePaymentWebhook processes payment webhook
func (h *CreditPaymentHandler) HandlePaymentWebhook(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	paymentID, err := h.parseID(c.Param("id"))
	if err != nil {
		h.logger.Error("Invalid payment ID", zap.Error(err))
		return h.handleError(c, err, "invalid payment ID")
	}

	var req struct {
		Status string `json:"status" validate:"required"`
	}
	if err := c.Bind(&req); err != nil {
		h.logger.Error("Failed to bind request body", zap.Error(err))
		return h.handleError(c, err, "invalid request body")
	}

	h.logger.Info("PaymentWebhook received", zap.Uint("paymentID", paymentID), zap.Any("request", req))

	if err := h.validator.Validate(req); err != nil {
		h.logger.Error("Validation failed", zap.Error(err))
		return h.handleError(c, err, "validation failed")
	}

	if err := h.service.HandlePaymentWebhook(ctx, paymentID, req.Status); err != nil {
		h.logger.Error("Failed to process payment webhook", zap.Error(err))
		return h.handleError(c, err, "failed to process payment webhook")
	}

	h.logger.Info("Payment webhook processed successfully", zap.Uint("paymentID", paymentID))
	return c.JSON(http.StatusOK, map[string]string{"message": "Payment webhook processed successfully"})
}
// HandleInstallmentWebhook processes installment webhook
func (h *CreditPaymentHandler) HandleInstallmentWebhook(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	installmentID, err := h.parseID(c.Param("id"))
	if err != nil {
		h.logger.Error("Invalid installment ID", zap.Error(err))
		return h.handleError(c, err, "invalid installment ID")
	}

	var req struct {
		Status string `json:"status" validate:"required"`
	}
	if err := c.Bind(&req); err != nil {
		h.logger.Error("Failed to bind request body", zap.Error(err))
		return h.handleError(c, err, "invalid request body")
	}

	h.logger.Info("InstallmentWebhook received", zap.Uint("installmentID", installmentID), zap.Any("request", req))

	if err := h.validator.Validate(req); err != nil {
		h.logger.Error("Validation failed", zap.Error(err))
		return h.handleError(c, err, "validation failed")
	}

	if err := h.service.HandleInstallmentWebhook(ctx, installmentID, req.Status); err != nil {
		h.logger.Error("Failed to process installment webhook", zap.Error(err))
		return h.handleError(c, err, "failed to process installment webhook")
	}

	h.logger.Info("Installment webhook processed successfully", zap.Uint("installmentID", installmentID))
	return c.JSON(http.StatusOK, map[string]string{"message": "Installment webhook processed successfully"})
}

func (h *CreditPaymentHandler) handleError(c echo.Context, err error, message string) error {
	h.logger.Error(message, zap.Error(err))

	switch {
	case errors.As(err, &utils.ErrNotFound{}):
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	case errors.As(err, &utils.ErrDuplicateEntry{}):
		return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
	case errors.As(err, &utils.ErrDatabase{}):
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "An unexpected error occurred"})
	case errors.Is(err, services.ErrInsufficientCredit):
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Insufficient credit"})
	case errors.Is(err, services.ErrInvalidAmount):
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid amount"})
	case errors.Is(err, services.ErrPaymentFailed):
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Payment processing failed"})
	default:
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "An unexpected error occurred"})
	}
}

func (h *CreditPaymentHandler) parseID(param string) (uint, error) {
	id, err := strconv.ParseUint(param, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

func (h *CreditPaymentHandler) createCreditApplicationResponse(app *domains.CreditApplication) dto.CreditApplicationResponse {
	return dto.CreditApplicationResponse{
		ID:        app.ID,
		UserID:    app.UserID,
		Amount:    app.Amount,
		Currency:  app.Currency,
		Status:    app.Status,
		CreatedAt: app.CreatedAt,
	}
}

func (h *CreditPaymentHandler) createPaymentResponse(payment *domains.Payment) dto.PaymentResponse {
	resp := dto.PaymentResponse{
		ID:                  payment.ID,
		CreditApplicationID: payment.CreditApplicationID,
		UserID:              payment.UserID,
		Amount:              payment.Amount,
		Currency:            payment.Currency,
		PaymentMethod:       payment.PaymentMethod,
		Status:              payment.Status,
		CreatedAt:           payment.CreatedAt,
		Installments:        make([]dto.InstallmentResponse, len(payment.Installments)),
	}

	for i, installment := range payment.Installments {
		resp.Installments[i] = dto.InstallmentResponse{
			ID:                installment.ID,
			PaymentID:         installment.PaymentID,
			InstallmentNumber: installment.InstallmentNumber,
			DueDate:           installment.DueDate,
			Amount:            installment.Amount,
			Status:            installment.Status,
			CreatedAt:         installment.CreatedAt,
		}
	}

	return resp
}