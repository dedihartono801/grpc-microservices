package grpc

import (
	"context"

	"github.com/dedihartono801/product-svc/internal/app/usecase/product"
	"github.com/dedihartono801/product-svc/pkg/customstatus"
	"github.com/dedihartono801/product-svc/pkg/dto"
	pb "github.com/dedihartono801/product-svc/pkg/protobuf"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ProductHandler struct {
	Service product.Service
}

func (h *ProductHandler) ListProduct(ctx context.Context, in *emptypb.Empty) (*pb.ProductResponse, error) {
	products, err := h.Service.ListProduct()
	if err != nil {
		return &pb.ProductResponse{
			Status: int32(customstatus.ErrInternalServerError.Code),
			Error:  err.Error(),
		}, nil
	}

	return &pb.ProductResponse{
		Status: int32(customstatus.StatusOk.Code),
		Data:   products,
	}, nil
}

func (h *ProductHandler) ProcessProduct(ctx context.Context, in *pb.ProcessProductRequest) (*pb.ProcessProductResponse, error) {
	req := &dto.ProcessProductRequestDto{
		ProductId: in.ProductId,
		Quantity:  in.Quantity,
	}
	result, statusCode, err := h.Service.ProcessProduct(req)
	if err != nil {
		return &pb.ProcessProductResponse{
			Status: int32(statusCode),
			Error:  err.Error(),
		}, nil
	}

	return &pb.ProcessProductResponse{
		Status: int32(statusCode),
		Price:  result.Price,
	}, nil
}
