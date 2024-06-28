package domain

import (
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UniversalId           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()" json:"universal_id"`
	FirstName             string    `db:"first_name" json:"first_name"`
	LastName              string    `db:"last_name" json:"last_name"`
	Email                 string    `db:"email" json:"email"`
	PasswordHash          string    `db:"password_hash" json:"-"`
	PhoneNumber           string    `db:"phone_number" json:"phone_number"`
	Address               string    `db:"address" json:"address"`
	LoyaltyPoints         int       `db:"loyalty_points" json:"loyalty_points"`
	RefreshToken          string    `db:"refresh_token" json:"refresh_token"`
	RefreshTokenExpiresAt time.Time `db:"refresh_token_exp_date" json:"refresh_token_exp_date"`
	IDImageURL            string    `db:"id_image_url" json:"id_image_url"`
}
