package repositories

import (
	"github.com/gofrs/uuid"
	domain "github.com/mohamed2394/sahla/internal/domains"

)

type UserRepository interface {
	Create(user *domain.User) error
	GetByID(id uuid.UUID) (*domain.User, error)
	GetByEmail(email string) (*domain.User, error)
	Update(user *domain.User) error
	Delete(id uuid.UUID) error
	List(offset, limit int) ([]*domain.User, error)
	FindByCriteria(criteria map[string]interface{}) ([]*domain.User, error)
	UpdateIDImage(id uuid.UUID, imageURL string) error

	// TODO
}
