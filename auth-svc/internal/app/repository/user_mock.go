package repository

import (
	"errors"

	"github.com/dedihartono801/auth-svc/internal/entity"
	"github.com/dedihartono801/auth-svc/pkg/helpers"
)

// Define the MockAdminRepository interface
type MockUserRepository interface {
	Create(transaction *entity.User) error
	FindByEmail(email string) (*entity.User, error)
}

type mockUserRepository struct{}

func NewMockUserRepository() MockUserRepository {
	return &mockUserRepository{}

}

func (m *mockUserRepository) FindByEmail(email string) (*entity.User, error) {
	account := &entity.User{
		ID:       1,
		Username: "diding",
		Email:    "diding@gmail.com",
		Password: helpers.EncryptPassword("123"),
	}
	if account.Email == email {
		return account, nil
	}
	return nil, errors.New("account not found")
}

func (m *mockUserRepository) Create(account *entity.User) error {
	return nil
}
