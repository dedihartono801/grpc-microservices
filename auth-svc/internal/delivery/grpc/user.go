package grpc

import (
	"context"

	"github.com/dedihartono801/auth-svc/internal/app/usecase/user"
	"github.com/dedihartono801/auth-svc/pkg/dto"
	pb "github.com/dedihartono801/auth-svc/pkg/protobuf"
)

type UserHandler struct {
	Service user.Service
}

func (h *UserHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	userDto := &dto.UserCreateRequestDto{
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
	}
	account, statuscode, err := h.Service.Register(userDto)
	if err != nil {
		return &pb.RegisterResponse{
			Status: int32(statuscode),
			Error:  err.Error(),
		}, nil
	}

	data := &pb.User{
		Username: account.Username,
		Email:    account.Email,
	}
	return &pb.RegisterResponse{
		Status: int32(statuscode),
		Data:   data,
	}, nil
}

func (h *UserHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	loginDto := &dto.UserLoginRequestDto{
		Email:    req.Email,
		Password: req.Password,
	}

	result, statusCode, err := h.Service.Login(loginDto)
	if err != nil {
		return &pb.LoginResponse{
			Status: int32(statusCode),
			Error:  err.Error(),
		}, nil
	}
	data := &pb.Token{
		Token:     result.Token,
		ExpiredAt: result.ExpiredAt,
	}
	return &pb.LoginResponse{
		Status: int32(statusCode),
		Data:   data,
	}, nil

}

func (h *UserHandler) Validate(ctx context.Context, req *pb.ValidateRequest) (*pb.ValidateResponse, error) {

	result, statusCode, err := h.Service.Validate(req.Token)
	if err != nil {
		return &pb.ValidateResponse{
			Status: int32(statusCode),
			Error:  err.Error(),
		}, nil
	}
	return &pb.ValidateResponse{
		Status: int32(statusCode),
		UserId: result.UserId,
	}, nil

}
