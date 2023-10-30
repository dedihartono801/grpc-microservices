package transaction

import (
	"context"
	"errors"
	"sync"

	productSvc "github.com/dedihartono801/transaction-svc/cmd/grpc/client/product"
	productPb "github.com/dedihartono801/transaction-svc/cmd/grpc/client/product/pb"
	"github.com/dedihartono801/transaction-svc/internal/app/repository"
	"github.com/dedihartono801/transaction-svc/internal/entity"
	"github.com/dedihartono801/transaction-svc/pkg/customstatus"
	"github.com/dedihartono801/transaction-svc/pkg/dto"
	"github.com/dedihartono801/transaction-svc/pkg/identifier"
	pb "github.com/dedihartono801/transaction-svc/pkg/protobuf"
	"github.com/dedihartono801/transaction-svc/pkg/validator"
	"golang.org/x/sync/errgroup"
)

type Service interface {
	CreateTransaction(input *dto.TransactionRequestDto) (int64, int, error)
	GetDetailTransaction(input *pb.GetDetailTransactionRequest) (*pb.Transaction, int, error)
}

type service struct {
	transactionRepository   repository.TransactionRepository
	dbTransactionRepository repository.DbTransactionRepository
	productSvcClient        productSvc.ServiceClient
	validator               validator.Validator
	identifier              identifier.Identifier
}

func NewGrpcTransactionService(
	transactionRepository repository.TransactionRepository,
	dbTransactionRepository repository.DbTransactionRepository,
	productSvcClient productSvc.ServiceClient,
	validator validator.Validator,
	identifier identifier.Identifier,
) Service {
	return &service{
		transactionRepository:   transactionRepository,
		dbTransactionRepository: dbTransactionRepository,
		productSvcClient:        productSvcClient,
		validator:               validator,
		identifier:              identifier,
	}
}

func (s *service) CreateTransaction(input *dto.TransactionRequestDto) (int64, int, error) {
	if err := s.validator.Validate(input); err != nil {
		return 0, customstatus.ErrBadRequest.Code, err
	}

	// Create the transaction
	transaction := &entity.Transaction{
		UserId: input.UserId,
	}
	// Begin a database transaction
	tx, err := s.dbTransactionRepository.BeginTransaction()
	if err != nil {
		return 0, customstatus.ErrInternalServerError.Code, errors.New(customstatus.ErrInternalServerError.Message)
	}

	// Create the transaction and transaction items
	trxId, err := s.transactionRepository.CreateTransaction(tx, transaction)
	if err != nil {
		return 0, customstatus.ErrInternalServerError.Code, errors.New(customstatus.ErrInternalServerError.Message)
	}

	dt := make([]*entity.TransactionItem, len(input.Items))
	totalPrice := 0
	totalQuantity := 0
	var mutex sync.Mutex
	ctx := context.Background()
	// create an errgroup.Group instance
	var g errgroup.Group
	for i, items := range input.Items {
		item := items
		i := i

		//process item using goroutine
		g.Go(func() error {

			// Check stock product in product service
			product, err := s.productSvcClient.Product.ProcessProduct(ctx, &productPb.ProcessProductRequest{
				ProductId: item.ProductId,
				Quantity:  item.Quantity,
			})

			if err != nil {
				return errors.New(err.Error())
			}

			if product.Error != "" {
				return errors.New(product.Error)
			}

			dt[i] = &entity.TransactionItem{
				TransactionId: trxId,
				ProductId:     item.ProductId,
				Price:         product.Price,
				Quantity:      item.Quantity,
			}

			//mutex to avoid race condition
			mutex.Lock()
			totalPrice += int(product.Price) * int(item.Quantity)
			mutex.Unlock()

			return nil
		})
		totalQuantity += int(item.Quantity)

	}

	// wait for all goroutines to finish
	if err := g.Wait(); err != nil {
		tx.Rollback()
		return 0, customstatus.ErrBadRequest.Code, errors.New(err.Error())
	}

	if err := s.transactionRepository.CreateTransactionItem(tx, dt); err != nil {
		tx.Rollback()
		return 0, customstatus.ErrInternalServerError.Code, errors.New(err.Error())
	}

	transactionDt, err := s.transactionRepository.GetTransactionById(tx, trxId)
	if err != nil {
		tx.Rollback()
		return 0, customstatus.ErrInternalServerError.Code, errors.New(customstatus.ErrInternalServerError.Message)
	}

	transactionDt.TotalAmount = int64(totalPrice)
	transactionDt.TotalQuantity = int64(totalQuantity)
	err = s.transactionRepository.UpdateTransaction(tx, transactionDt)
	if err != nil {
		tx.Rollback()
		return 0, customstatus.ErrInternalServerError.Code, errors.New(customstatus.ErrInternalServerError.Message)
	}

	// Commit the database transaction
	if err := s.dbTransactionRepository.CommitTransaction(tx); err != nil {
		return 0, customstatus.ErrInternalServerError.Code, errors.New(customstatus.ErrInternalServerError.Message)
	}
	return trxId, int(customstatus.StatusCreated.Code), nil
}

func (s *service) GetDetailTransaction(input *pb.GetDetailTransactionRequest) (*pb.Transaction, int, error) {
	transaction, err := s.transactionRepository.GetTransaction(input.TransactionId)
	if err != nil {
		return nil, customstatus.ErrNotFound.Code, errors.New(customstatus.ErrNotFound.Message)
	}

	transactionItems, err := s.transactionRepository.GetTransactionItemByTransactionId(input.TransactionId)
	if err != nil {
		return nil, customstatus.ErrNotFound.Code, errors.New(customstatus.ErrNotFound.Message)
	}

	data := &pb.Transaction{
		Id:            transaction.ID,
		UserId:        transaction.UserId,
		TotalAmount:   transaction.TotalAmount,
		TotalQuantity: transaction.TotalQuantity,
		Items:         transactionItems,
	}

	return data, customstatus.StatusOk.Code, nil
}
