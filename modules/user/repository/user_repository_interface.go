package repository

import (
	"github.com/mohamed2394/sahla/modules/user/domain"
)

type UserRepository interface {
	Create(user *domain.User) error
	GetByID(id uint) (*domain.User, error)
	GetByEmail(email string) (*domain.User, error)
	Update(user *domain.User) error
	Delete(id uint) error
	List(offset, limit int) ([]*domain.User, error)
	FindByCriteria(criteria map[string]interface{}) ([]*domain.User, error)

	// TODO
	//CreateUserWithInitialLoan(user *domain.User, loan *domain.Loan) error
}
