package user

import (
	"errors"
	"time"

	"github.com/dedihartono801/auth-svc/internal/app/repository"
	"github.com/dedihartono801/auth-svc/internal/entity"
	"github.com/dedihartono801/auth-svc/pkg/config"
	"github.com/dedihartono801/auth-svc/pkg/customstatus"
	"github.com/dedihartono801/auth-svc/pkg/dto"
	"github.com/dedihartono801/auth-svc/pkg/helpers"
	"github.com/dedihartono801/auth-svc/pkg/identifier"
	"github.com/dedihartono801/auth-svc/pkg/validator"
	"github.com/golang-jwt/jwt"
)

type Service interface {
	Register(input *dto.UserCreateRequestDto) (*entity.User, int, error)
	Login(input *dto.UserLoginRequestDto) (*dto.LoginResponse, int, error)
	Validate(bearerToken string) (*dto.ValidateResponse, int, error)
}

type service struct {
	userRepository repository.UserRepository
	validator      validator.Validator
	identifier     identifier.Identifier
}

func NewGrpcAccountService(
	userRepository repository.UserRepository,
	validator validator.Validator,
	identifier identifier.Identifier,
) Service {
	return &service{
		userRepository: userRepository,
		validator:      validator,
		identifier:     identifier,
	}
}

func (s *service) Register(input *dto.UserCreateRequestDto) (*entity.User, int, error) {
	if err := s.validator.Validate(input); err != nil {
		return nil, customstatus.ErrBadRequest.Code, errors.New(customstatus.ErrBadRequest.Message)
	}
	ok := helpers.ValidateInputSpecific(input.Email)
	if !ok {
		return nil, customstatus.ErrBadRequest.Code, errors.New("invalid input email")
	}

	ok = helpers.ValidateInputCommon(input.Username)
	if !ok {
		return nil, customstatus.ErrBadRequest.Code, errors.New("invalid input username")
	}

	ok = helpers.ValidateInputSpecific(input.Password)
	if !ok {
		return nil, customstatus.ErrBadRequest.Code, errors.New("invalid input password")
	}

	if input.Password != "" {
		input.Password = helpers.EncryptPassword(input.Password)
	}
	user := entity.User{
		Username: input.Username,
		Email:    input.Email,
		Password: input.Password,
	}

	_, err := s.userRepository.FindByEmail(user.Email)
	if err == nil {
		return nil, customstatus.ErrEmailFound.Code, errors.New(customstatus.ErrEmailFound.Message)
	}

	err = s.userRepository.Create(&user)
	if err != nil {
		return nil, customstatus.ErrInternalServerError.Code, errors.New(customstatus.ErrInternalServerError.Message)
	}
	return &user, customstatus.StatusCreated.Code, nil
}

func (s *service) Login(input *dto.UserLoginRequestDto) (*dto.LoginResponse, int, error) {
	if err := s.validator.Validate(input); err != nil {
		return nil, customstatus.ErrBadRequest.Code, errors.New(customstatus.ErrBadRequest.Message)
	}
	ok := helpers.ValidateInputSpecific(input.Email)
	if !ok {
		return nil, customstatus.ErrBadRequest.Code, errors.New("invalid input email")
	}

	user, err := s.userRepository.FindByEmail(input.Email)
	if err != nil {
		return nil, customstatus.ErrEmailNotFound.Code, errors.New(customstatus.ErrEmailNotFound.Message)
	}

	if user.Password != helpers.EncryptPassword(input.Password) {
		return nil, customstatus.ErrPasswordWrong.Code, errors.New(customstatus.ErrPasswordWrong.Message)
	}
	expirationTime := time.Now().Add(time.Hour * time.Duration(24))
	token, err := helpers.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, customstatus.ErrInternalServerError.Code, errors.New(customstatus.ErrInternalServerError.Message)
	}
	responseParams := dto.LoginResponse{
		Token:     token,
		ExpiredAt: expirationTime.Format(time.RFC3339),
	}
	return &responseParams, customstatus.StatusOk.Code, nil
}

func (s *service) Validate(bearerToken string) (*dto.ValidateResponse, int, error) {

	// Verify that the Authorization header starts with "Bearer "
	if len(bearerToken) < 7 || bearerToken[:7] != "Bearer " {
		return nil, customstatus.ErrUnAuthorized.Code, errors.New("invalid format authorization")
	}

	// Parse the JWT from the Authorization header
	tokenString := bearerToken[len("Bearer "):]
	token, err := jwt.ParseWithClaims(tokenString, &helpers.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GetEnv("JWT_SECRET")), nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, customstatus.ErrUnAuthorized.Code, errors.New("invalid signature")
		}
		return nil, customstatus.ErrUnAuthorized.Code, errors.New(customstatus.ErrUnAuthorized.Message)
	}

	// Check if the JWT has expired
	claims, ok := token.Claims.(*helpers.Claims)
	if !ok || !token.Valid {
		return nil, customstatus.ErrUnAuthorized.Code, errors.New(customstatus.ErrUnAuthorized.Message)
	}
	if claims.ExpiresAt < time.Now().Unix() {
		return nil, customstatus.ErrUnAuthorized.Code, errors.New(customstatus.ErrUnAuthorized.Message)
	}

	return &dto.ValidateResponse{
		UserId: claims.UserId,
	}, customstatus.StatusOk.Code, nil
}
