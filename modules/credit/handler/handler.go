package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	domain "github.com/mohamed2394/sahla/modules/credit"
	"github.com/mohamed2394/sahla/modules/credit/service"
)

type CreditHandler struct {
	creditService *service.CreditService
}

func NewCreditHandler(creditService *service.CreditService) *CreditHandler {
	return &CreditHandler{creditService: creditService}
}

func (h *CreditHandler) SaveCreditScore(c echo.Context) error {
	var score domain.CreditScore
	if err := c.Bind(&score); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := h.creditService.SaveCreditScore(&score); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, score)
}

func (h *CreditHandler) GetCreditScoreByUserID(c echo.Context) error {
	userID := c.Param("userID")
	score, err := h.creditService.GetCreditScoreByUserID(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	if score == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Credit score not found"})
	}

	return c.JSON(http.StatusOK, score)
}

func (h *CreditHandler) SaveCreditFeatures(c echo.Context) error {
	var features domain.CreditFeatures
	if err := c.Bind(&features); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	creditScore, err := h.creditService.SaveCreditFeatures(&features)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, creditScore)
}

func (h *CreditHandler) GetCreditFeaturesByUserID(c echo.Context) error {
	userID := c.Param("userID")
	features, err := h.creditService.GetCreditFeaturesByUserID(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	if features == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Credit features not found"})
	}

	return c.JSON(http.StatusOK, features)
}
