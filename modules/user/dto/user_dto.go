package dto

type CreateUserRequest struct {
	FirstName     string `json:"first_name" validate:"required"`
	LastName      string `json:"last_name" validate:"required"`
	Email         string `json:"email" validate:"required,email"`
	PhoneNumber   string `json:"phone_number"`
	Address       string `json:"address"`
	Password      string `json:"password" validate:"required"`
	LoyaltyPoints int    `json:"loyalty_points"`
}

type UpdateUserRequest struct {
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Password      string `json:"password"`
	Email         string `json:"email"`
	PhoneNumber   string `json:"phone_number"`
	Address       string `json:"address"`
	LoyaltyPoints int    `json:"loyalty_points"`
}
type UserResponse struct {
	ID             uint   `json:"id"`
	FirstName      string `json:"first_name" validate:"required"`
	LastName       string `json:"last_name" validate:"required"`
	Email          string `json:"email"`
	PhoneNumber    string `json:"phone_number"`
	Address        string `json:"address"`
	LoyaltyPoints  int    `json:"loyalty_points"`
	MembershipType string `json:"membership_type"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}
