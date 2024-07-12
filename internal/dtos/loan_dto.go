package dto

import (
	"time"

	domains "github.com/mohamed2394/sahla/internal/domains"

	"github.com/gofrs/uuid"
)

type CreateLoanRequest struct {
	UserID uint    `json:"user_id" validate:"required"`
	Amount float64 `json:"amount" validate:"required,gt=0"`
	Term   int     `json:"term" validate:"required,gt=0"`
}

type LoanResponse struct {
	ID                uint              `json:"id"`
	UniversalID       uuid.UUID         `json:"universal_id"`
	UserID            uint              `json:"user_id"`
	Amount            float64           `json:"amount"`
	OriginalTerm      int               `json:"original_term"`
	RemainingTerm     int               `json:"remaining_term"`
	Status            domains.LoanStatus `json:"status"`
	ApprovedAt        *time.Time        `json:"approved_at,omitempty"`
	ActivatedAt       *time.Time        `json:"activated_at,omitempty"`
	NextPaymentDate   *time.Time        `json:"next_payment_date,omitempty"`
	RemainingAmount   float64           `json:"remaining_amount"`
	TotalPaidAmount   float64           `json:"total_paid_amount"`
	CreditScoreImpact int               `json:"credit_score_impact"`
}

type LoanDetailResponse struct {
	LoanResponse
	Payments          []PaymentResponse       `json:"payments"`
	RepaymentSchedule []RepaymentPlanResponse `json:"repayment_schedule"`
}
