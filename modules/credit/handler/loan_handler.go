// internal/handler/loan_handler.go
package handler

import (
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/mohamed2394/sahla/modules/credit/dto"
	"github.com/mohamed2394/sahla/modules/credit/service"
)

type LoanHandler struct {
	loanService service.LoanService
}

func NewLoanHandler(loanService service.LoanService) *LoanHandler {
	return &LoanHandler{
		loanService: loanService,
	}
}

func (h *LoanHandler) GetLoansByUserID(c echo.Context) error {
	// Get the user claims from the context
	claims, ok := c.Get("user").(jwt.MapClaims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token claims"})
	}

	// Get the user ID from the claims
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid user ID in token"})
	}

	loans, err := h.loanService.GetLoansByUserID(c.Request().Context(), uint(userID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, loans)
}

func (h *LoanHandler) CreateLoan(c echo.Context) error {
	var req dto.CreateLoanRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	loan, err := h.loanService.CreateLoan(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, loan)
}

func (h *LoanHandler) GetLoanByID(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid loan ID"})
	}

	loan, err := h.loanService.GetLoanByID(c.Request().Context(), uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, loan)
}

func (h *LoanHandler) GetLoanByUniversalID(c echo.Context) error {
	universalID := c.Param("universalID")

	loan, err := h.loanService.GetLoanByUniversalID(c.Request().Context(), universalID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, loan)
}

func (h *LoanHandler) UpdateLoanStatus(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid loan ID"})
	}

	var req struct {
		Status string `json:"status"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	err = h.loanService.UpdateLoanStatus(c.Request().Context(), uint(id), req.Status)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusOK)
}
