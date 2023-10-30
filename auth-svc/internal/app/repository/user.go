package repository

import (
	"github.com/dedihartono801/auth-svc/internal/entity"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(transaction *entity.User) error
	FindByEmail(email string) (*entity.User, error)
}

type userRepository struct {
	database *gorm.DB
}

func NewUserRepository(database *gorm.DB) UserRepository {
	return &userRepository{database}
}

func (r *userRepository) Create(user *entity.User) error {
	return r.database.Table("users").Create(user).Error
}

func (r *userRepository) FindByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := r.database.Where("email = ?", email).First(&user).Error
	return &user, err
}
