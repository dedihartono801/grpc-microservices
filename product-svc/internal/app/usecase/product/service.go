package product

import (
	"errors"
	"fmt"

	"github.com/dedihartono801/product-svc/internal/app/repository"
	"github.com/dedihartono801/product-svc/pkg/customstatus"
	"github.com/dedihartono801/product-svc/pkg/dto"
	"github.com/dedihartono801/product-svc/pkg/identifier"
	pb "github.com/dedihartono801/product-svc/pkg/protobuf"
	"github.com/dedihartono801/product-svc/pkg/validator"
)

type Service interface {
	ListProduct() ([]*pb.Product, error)
	ProcessProduct(input *dto.ProcessProductRequestDto) (*pb.ProcessProductResponse, int, error)
}

type service struct {
	productRepository repository.ProductRepository
	validator         validator.Validator
	identifier        identifier.Identifier
}

func NewGrpcProductService(
	productRepository repository.ProductRepository,
	validator validator.Validator,
	identifier identifier.Identifier,
) Service {
	return &service{
		productRepository: productRepository,
		validator:         validator,
		identifier:        identifier,
	}
}

func (s *service) ListProduct() ([]*pb.Product, error) {
	var products []*pb.Product
	result, err := s.productRepository.ListProduct()
	if err != nil {
		return nil, errors.New(customstatus.ErrInternalServerError.Message)
	}

	for _, val := range result {
		products = append(products, &pb.Product{
			Id:    val.ID,
			Name:  val.Name,
			Stock: val.Stock,
			Price: val.Price,
		})
	}

	return products, nil
}

func (s *service) ProcessProduct(input *dto.ProcessProductRequestDto) (*pb.ProcessProductResponse, int, error) {
	if err := s.validator.Validate(input); err != nil {
		return nil, customstatus.ErrBadRequest.Code, errors.New(customstatus.ErrBadRequest.Message)
	}
	product, err := s.productRepository.GetProductById(input.ProductId)
	if err != nil {
		return nil, customstatus.ErrNotFound.Code, errors.New(customstatus.ErrNotFound.Message)
	}

	if product.Stock < input.Quantity {
		err := fmt.Errorf("Stock tidak mencukupi untuk product %s sisa stock adalah %d", product.Name, product.Stock)
		return nil, customstatus.ErrBadRequest.Code, errors.New(err.Error())
	}

	product.Stock = product.Stock - input.Quantity
	err = s.productRepository.UpdateStockProduct(product)
	if err != nil {
		return nil, customstatus.ErrInternalServerError.Code, errors.New(customstatus.ErrInternalServerError.Message)
	}

	return &pb.ProcessProductResponse{
		Price: product.Price,
	}, customstatus.StatusCreated.Code, nil
}
