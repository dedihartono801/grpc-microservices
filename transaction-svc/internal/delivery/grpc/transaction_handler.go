package grpc

import (
	"context"
	"errors"

	"github.com/dedihartono801/transaction-svc/internal/app/usecase/transaction"
	"github.com/dedihartono801/transaction-svc/pkg/dto"
	pb "github.com/dedihartono801/transaction-svc/pkg/protobuf"
)

type TransactionHandler struct {
	Service transaction.Service
}

func (h *TransactionHandler) Transaction(ctx context.Context, input *pb.TransactionRequest) (*pb.TransactionResponse, error) {
	req := &dto.TransactionRequestDto{
		UserId: input.UserId,
		Items:  input.Items,
	}
	trxId, statusCode, err := h.Service.CreateTransaction(req)
	if err != nil {
		return &pb.TransactionResponse{
			Status: int32(statusCode),
			Error:  err.Error(),
		}, errors.New(err.Error())
	}

	return &pb.TransactionResponse{
		Status:        int32(statusCode),
		TransactionId: trxId,
	}, nil

}

func (h *TransactionHandler) GetDetailTransaction(ctx context.Context, req *pb.GetDetailTransactionRequest) (*pb.DetailTransactionResponse, error) {

	dt, statusCode, err := h.Service.GetDetailTransaction(req)
	if err != nil {
		return &pb.DetailTransactionResponse{
			Status: int32(statusCode),
			Error:  err.Error(),
		}, errors.New(err.Error())
	}

	return &pb.DetailTransactionResponse{
		Status: int32(statusCode),
		Data:   dt,
	}, nil

}
