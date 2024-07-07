package domain

import (
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "PENDING"
	PaymentStatusCompleted PaymentStatus = "COMPLETED"
	PaymentStatusFailed    PaymentStatus = "FAILED"
	PaymentStatusRefunded  PaymentStatus = "REFUNDED"
)

type Payment struct {
	gorm.Model
	UniversalID     uuid.UUID     `gorm:"type:uuid;default:gen_random_uuid()" json:"universal_id"`
	LoanID          uint          `json:"loan_id"`
	Amount          float64       `json:"amount"`
	PaymentDate     time.Time     `json:"payment_date"`
	Status          PaymentStatus `json:"status"`
	PaymentMethod   string        `json:"payment_method"`
	TransactionID   string        `json:"transaction_id"`
	ProcessingFee   float64       `json:"processing_fee"`
	LateFee         float64       `json:"late_fee"`
	EarlyPaymentFee float64       `json:"early_payment_fee"`
}

func (p *Payment) BeforeCreate(tx *gorm.DB) (err error) {
	if p.UniversalID == uuid.Nil {
		p.UniversalID, err = uuid.NewV4()
	}
	return
}
