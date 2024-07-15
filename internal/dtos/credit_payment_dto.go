package dtos

import (
	"time"
	"github.com/mohamed2394/sahla/internal/domains"
)

// CreditApplicationRequest represents the DTO for creating a credit application
type CreditApplicationRequest struct {
	UserID   string `json:"user_id" binding:"required"`
	Amount   int    `json:"amount" binding:"required,min=1"`
	Currency string `json:"currency" binding:"required,len=3"`
}

// CreditApplicationResponse represents the DTO for credit application response
type CreditApplicationResponse struct {
	ID        uint      `json:"id"`
	UserID    string    `json:"user_id"`
	Amount    int       `json:"amount"`
	Currency  string    `json:"currency"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// PaymentRequest represents the DTO for creating a payment
type PaymentRequest struct {
	CreditApplicationID uint                  `json:"credit_application_id" binding:"required"`
	UserID              string                `json:"user_id" binding:"required"`
	Amount              int                   `json:"amount" binding:"required,min=1"`
	Currency            string                `json:"currency" binding:"required,len=3"`
	PaymentMethod       domains.PaymentMethod `json:"payment_method" binding:"required"`
}

// PaymentResponse represents the DTO for payment response
type PaymentResponse struct {
	ID                  uint                    `json:"id"`
	CreditApplicationID uint                    `json:"credit_application_id"`
	UserID              string                  `json:"user_id"`
	Amount              int                     `json:"amount"`
	Currency            string                  `json:"currency"`
	PaymentMethod       domains.PaymentMethod   `json:"payment_method"`
	Status              string                  `json:"status"`
	Installments        []InstallmentResponse   `json:"installments,omitempty"`
	CreatedAt           time.Time               `json:"created_at"`
}

// InstallmentResponse represents the DTO for installment response
type InstallmentResponse struct {
	ID                uint      `json:"id"`
	PaymentID         uint      `json:"payment_id"`
	InstallmentNumber int       `json:"installment_number"`
	DueDate           string    `json:"due_date"`
	Amount            int       `json:"amount"`
	Status            string    `json:"status"`
	CreatedAt         time.Time `json:"created_at"`
}