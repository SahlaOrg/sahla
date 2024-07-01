package domain

import (
	"time"
)

type CreditScore struct {
	UserID       string `gorm:"primaryKey"`
	Score        int
	ScoreDate    time.Time
	RiskCategory string
	ScoreFactors []string `gorm:"type:text[]"`
}

type CreditAssessment struct {
	UserID           string `gorm:"primaryKey"`
	AssessmentDate   time.Time
	CriteriaMet      map[string]bool `gorm:"type:jsonb"`
	AssessmentResult string
	RequestedAmount  float64
	ApprovedAmount   float64
	ExpirationDate   time.Time
}

type CreditApplication struct {
	ApplicationID   string `gorm:"primaryKey"`
	UserID          string
	RequestedAmount float64
	Purpose         string
	ApplicationDate time.Time
	Status          string
}

type CreditLimit struct {
	UserID          string `gorm:"primaryKey"`
	CurrentLimit    float64
	AvailableCredit float64
	LastUpdateDate  time.Time
}

type CreditFeatures struct {
	UserID              string  `gorm:"primaryKey"`
	IncomeLevel         float64 `json:"income_level"`
	DebtLevel           float64 `json:"debt_level"`
	CreditUtilization   float64 `json:"credit_utilization"`
	CreditHistoryLength int     `json:"credit_history_length"`
	NumCreditAccounts   int     `json:"num_credit_accounts"`
	NumCreditInquiries  int     `json:"num_credit_inquiries"`
	Age                 int     `json:"age"`
	PaymentHistory      string  `json:"payment_history"`
	EmploymentStatus    string  `json:"employment_status"`
	EducationLevel      string  `json:"education_level"`
}
