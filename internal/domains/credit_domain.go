package domains

import (
	"gorm.io/gorm"
)

// CreditApplication represents a credit application in the BNPL system.
type CreditApplication struct {
	gorm.Model
	UserID   string    `gorm:"type:uuid;not null" json:"user_id"`
	Amount   int       `gorm:"not null" json:"amount"`
	Currency string    `gorm:"type:varchar(3);not null" json:"currency"`
	Status   string    `gorm:"type:varchar(20);not null" json:"status"`
	Payments []Payment `json:"payments"`
}

// Payment represents a payment made towards a credit application.
type Payment struct {
	gorm.Model
	CreditApplicationID uint   `gorm:"not null" json:"credit_application_id"`
	UserID              string `gorm:"type:uuid;not null" json:"user_id"`
	OrderID             string `gorm:"type:uuid;not null;unique" json:"order_id"`
	Amount              int    `gorm:"not null" json:"amount"`
	Currency            string `gorm:"type:varchar(3);not null" json:"currency"`
	PaymentMethod       PaymentMethod `gorm:"embedded" json:"payment_method"`
	Status              string        `gorm:"type:varchar(20);not null" json:"status"`
	Installments        []Installment `json:"installments"`
}

// PaymentMethod represents the payment method details.
type PaymentMethod struct {
	Type    string         `gorm:"type:varchar(20);not null" json:"type"`
	Details PaymentDetails `gorm:"embedded" json:"details"`
}

// PaymentDetails represents the details of the payment method.
type PaymentDetails struct {
	CardNumber     string `gorm:"type:varchar(16)" json:"card_number"`
	CardHolderName string `gorm:"type:varchar(100)" json:"card_holder_name"`
	ExpiryDate     string `gorm:"type:varchar(5)" json:"expiry_date"`
	CVV            string `gorm:"type:varchar(4)" json:"cvv"`
}

// Installment represents an installment in the payment plan.
type Installment struct {
	gorm.Model
	PaymentID         uint   `gorm:"not null" json:"payment_id"`
	InstallmentNumber int    `gorm:"not null" json:"installment_number"`
	DueDate           string `gorm:"type:date;not null" json:"due_date"`
	Amount            int    `gorm:"not null" json:"amount"`
	Status            string `gorm:"type:varchar(20);not null" json:"status"`
}