package domain

import (
	"time"

	"gorm.io/gorm"
)

type RepaymentPlan struct {
	gorm.Model
	LoanID         uint      `json:"loan_id"`
	InstallmentNum int       `json:"installment_num"`
	DueDate        time.Time `json:"due_date"`
	Amount         float64   `json:"amount"`
	PrincipalPart  float64   `json:"principal_part"`
	Status         string    `json:"status"`
}
