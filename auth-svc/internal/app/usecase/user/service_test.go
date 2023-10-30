package user

import (
	"testing"
	"time"

	repoMock "github.com/dedihartono801/auth-svc/internal/app/repository"
	"github.com/dedihartono801/auth-svc/internal/entity"
	"github.com/dedihartono801/auth-svc/pkg/dto"
	"github.com/dedihartono801/auth-svc/pkg/helpers"
	"github.com/dedihartono801/auth-svc/pkg/identifier"
	"github.com/dedihartono801/auth-svc/pkg/validator"
	validatorv10 "github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	repo := repoMock.NewMockUserRepository()
	identifier := identifier.NewIdentifier()
	validator := validator.NewValidator(validatorv10.New())
	srv := NewGrpcAccountService(repo, validator, identifier)

	request := &dto.UserCreateRequestDto{
		Username: "diding",
		Email:    "diding2@gmail.com",
		Password: "123",
	}

	requestFailed := &dto.UserCreateRequestDto{
		Username: "diding",
		Email:    "diding@gmail.com",
		Password: "123",
	}

	expected := &entity.User{
		Username: "diding",
		Email:    "diding2@gmail.com",
		Password: helpers.EncryptPassword("123"),
	}

	// Define test cases
	testCases := []struct {
		name       string
		request    *dto.UserCreateRequestDto
		expected   *entity.User
		wantErr    bool
		statusCode int
	}{
		{
			name:       "Success create account",
			request:    request,
			expected:   expected,
			wantErr:    false,
			statusCode: 201,
		},
		{
			name:       "Fail create account",
			request:    requestFailed,
			expected:   nil,
			wantErr:    true,
			statusCode: 409,
		},
	}

	// Run tests
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, statusCode, err := srv.Register(tc.request)
			assert.Equal(t, tc.expected, actual, "Expected and actual data should be equal")
			assert.Equal(t, tc.statusCode, statusCode, "Expected status code and actual status code should be equal")
			assert.Equal(t, tc.wantErr, err != nil, "Expected error and actual error should be equal")
		})
	}

}

func TestLogin(t *testing.T) {
	repo := repoMock.NewMockUserRepository()
	identifier := identifier.NewIdentifier()
	validator := validator.NewValidator(validatorv10.New())
	srv := NewGrpcAccountService(repo, validator, identifier)

	request := &dto.UserLoginRequestDto{
		Email:    "diding@gmail.com",
		Password: "123",
	}

	expirationTime := time.Now().Add(time.Hour * time.Duration(24))
	token, _ := helpers.GenerateToken(uint64(1), request.Email)
	expected := &dto.LoginResponse{
		Token:     token,
		ExpiredAt: expirationTime.Format(time.RFC3339),
	}
	//accountId := 1

	// Define test cases
	testCases := []struct {
		name       string
		request    *dto.UserLoginRequestDto
		expected   *dto.LoginResponse
		wantErr    bool
		statusCode int
	}{
		{
			name:       "Success login",
			request:    request,
			expected:   expected,
			wantErr:    false,
			statusCode: 200,
		},
	}

	// Run tests
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, statusCode, err := srv.Login(tc.request)
			assert.Equal(t, tc.expected, actual, "Expected and actual data should be equal")
			assert.Equal(t, tc.statusCode, statusCode, "Expected status code and actual status code should be equal")
			assert.Equal(t, tc.wantErr, err != nil, "Expected error and actual error should be equal")
		})
	}

}
