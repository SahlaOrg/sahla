package domain

import (
	"time"

	user "github.com/mohamed2394/sahla/modules/user/domain"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type LoanStatus string

const (
	LoanStatusPending   LoanStatus = "PENDING"
	LoanStatusApproved  LoanStatus = "APPROVED"
	LoanStatusRejected  LoanStatus = "REJECTED"
	LoanStatusActive    LoanStatus = "ACTIVE"
	LoanStatusPaid      LoanStatus = "PAID"
	LoanStatusDefaulted LoanStatus = "DEFAULTED"
)

type Loan struct {
	gorm.Model
	UniversalID       uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid()" json:"universal_id"`
	User              user.User  `json:"user_id"`
	Amount            float64    `json:"amount"`
	OriginalTerm      int        `json:"original_term"` // in months
	RemainingTerm     int        `json:"remaining_term"`
	Status            LoanStatus `json:"status"`
	ApprovedAt        *time.Time `json:"approved_at"`
	ActivatedAt       *time.Time `json:"activated_at"`
	NextPaymentDate   *time.Time `json:"next_payment_date"`
	RemainingAmount   float64    `json:"remaining_amount"`
	TotalPaidAmount   float64    `json:"total_paid_amount"`
	LatePaymentFee    float64    `json:"late_payment_fee"`
	EarlyPaymentFee   float64    `json:"early_payment_fee"`
	CreditScoreImpact int        `json:"credit_score_impact"`
	Payments          []Payment  `gorm:"foreignKey:LoanID" json:"payments"`
}

func (l *Loan) BeforeCreate(tx *gorm.DB) (err error) {
	if l.UniversalID == uuid.Nil {
		l.UniversalID, err = uuid.NewV4()
	}
	return
}
