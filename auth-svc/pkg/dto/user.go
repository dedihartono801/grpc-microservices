package dto

type UserCreateRequestDto struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserLoginRequestDto struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"Password" validate:"required"`
}

type LoginResponse struct {
	Token     string `json:"token"`
	ExpiredAt string `json:"expired_at"`
}

type ValidateResponse struct {
	UserId uint64 `json:"user_id"`
}
