package repository

import (
	"github.com/dedihartono801/transaction-svc/internal/entity"
	pb "github.com/dedihartono801/transaction-svc/pkg/protobuf"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	CreateTransaction(tx *gorm.DB, transaction *entity.Transaction) (int64, error)
	CreateTransactionItem(tx *gorm.DB, ti []*entity.TransactionItem) error
	UpdateTransaction(tx *gorm.DB, transaction *entity.Transaction) error
	GetTransactionById(tx *gorm.DB, trxId int64) (*entity.Transaction, error)
	GetTransactionItemByTransactionId(trxId int64) ([]*pb.ItemsDetail, error)
	GetTransaction(trxId int64) (*entity.Transaction, error)
}

type transactionRepository struct {
	database *gorm.DB
}

func NewTransactionRepository(database *gorm.DB) TransactionRepository {
	return &transactionRepository{database}
}

func (r *transactionRepository) CreateTransaction(tx *gorm.DB, transaction *entity.Transaction) (int64, error) {
	result := tx.Table("transactions").Create(transaction)
	if result.Error != nil {
		return 0, result.Error
	}

	// Fetch the ID of the inserted record from the database
	insertedID := transaction.ID // Assuming you have an ID field in your user struct

	return insertedID, nil
}

func (r *transactionRepository) CreateTransactionItem(tx *gorm.DB, ti []*entity.TransactionItem) error {
	return tx.Table("transaction_items").Create(ti).Error
}

func (r *transactionRepository) UpdateTransaction(tx *gorm.DB, transaction *entity.Transaction) error {

	return tx.Table("transactions").Save(transaction).Error
}

func (r *transactionRepository) GetTransactionById(tx *gorm.DB, trxId int64) (*entity.Transaction, error) {
	var transaction entity.Transaction
	err := tx.Table("transactions").Where("id = ?", trxId).First(&transaction).Error
	return &transaction, err
}

func (r *transactionRepository) GetTransaction(trxId int64) (*entity.Transaction, error) {
	var transaction entity.Transaction
	err := r.database.Table("transactions").Where("id = ?", trxId).First(&transaction).Error
	return &transaction, err
}

func (r *transactionRepository) GetTransactionItemByTransactionId(trxId int64) ([]*pb.ItemsDetail, error) {
	var transactionItem []*pb.ItemsDetail
	result := r.database.Table("transaction_items").
		Select("id, product_id, quantity, price").
		Where("transaction_id = ?", trxId).
		Scan(&transactionItem)
	if result.Error != nil {
		return nil, result.Error
	}
	return transactionItem, nil
}
