package repository

import (
	"errors"

	"github.com/gofrs/uuid"
	"github.com/mohamed2394/sahla/modules/user/domain"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) GetByID(id uuid.UUID) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("universal_id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func (r *userRepository) GetByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(user *domain.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id uuid.UUID) error {
	result := r.db.Delete(&domain.User{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}

func (r *userRepository) List(offset, limit int) ([]*domain.User, error) {
	var users []*domain.User
	err := r.db.Offset(offset).Limit(limit).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) FindByCriteria(criteria map[string]interface{}) ([]*domain.User, error) {
	var users []*domain.User
	query := r.db
	for key, value := range criteria {
		query = query.Where(key+" = ?", value)
	}
	err := query.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}
