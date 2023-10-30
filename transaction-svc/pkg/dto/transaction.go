package dto

import pb "github.com/dedihartono801/transaction-svc/pkg/protobuf"

type TransactionRequestDto struct {
	UserId int64       `json:"user_id" validate:"required"`
	Items  []*pb.Items `json:"items" validate:"required"`
}
