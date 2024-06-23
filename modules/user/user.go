package user

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name          string `gorm:"not null"`
	Email         string `gorm:"unique;not null"`
	PhoneNumber   string
	Address       string
	PasswordHash  string `gorm:"not null"`
	LoyaltyPoints int    `gorm:"default:0"`
}
