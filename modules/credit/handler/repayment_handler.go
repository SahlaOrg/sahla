// internal/handler/repayment_plan_handler.go
package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/mohamed2394/sahla/modules/credit/service"
)

type RepaymentPlanHandler struct {
	repaymentPlanService service.RepaymentPlanService
}

func NewRepaymentPlanHandler(repaymentPlanService service.RepaymentPlanService) *RepaymentPlanHandler {
	return &RepaymentPlanHandler{
		repaymentPlanService: repaymentPlanService,
	}
}

func (h *RepaymentPlanHandler) GetRepaymentPlanByLoanID(c echo.Context) error {
	loanID, err := strconv.ParseUint(c.Param("loanID"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid loan ID"})
	}

	repaymentPlan, err := h.repaymentPlanService.GetRepaymentPlanByLoanID(c.Request().Context(), uint(loanID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, repaymentPlan)
}
