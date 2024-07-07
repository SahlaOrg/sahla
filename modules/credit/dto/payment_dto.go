package dto

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/mohamed2394/sahla/modules/credit/domain"
)

type CreatePaymentRequest struct {
	LoanID        uint    `json:"loan_id" validate:"required"`
	Amount        float64 `json:"amount" validate:"required,gt=0"`
	PaymentMethod string  `json:"payment_method" validate:"required"`
}

type PaymentResponse struct {
	ID              uint                 `json:"id"`
	UniversalID     uuid.UUID            `json:"universal_id"`
	LoanID          uint                 `json:"loan_id"`
	Amount          float64              `json:"amount"`
	PaymentDate     time.Time            `json:"payment_date"`
	Status          domain.PaymentStatus `json:"status"`
	PaymentMethod   string               `json:"payment_method"`
	TransactionID   string               `json:"transaction_id"`
	ProcessingFee   float64              `json:"processing_fee"`
	LateFee         float64              `json:"late_fee"`
	EarlyPaymentFee float64              `json:"early_payment_fee"`
}
