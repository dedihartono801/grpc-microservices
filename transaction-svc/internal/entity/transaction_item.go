package entity

import (
	"time"
)

type TransactionItem struct {
	ID            int64      `gorm:"primaryKey;autoIncrement:true" json:"id"`
	TransactionId int64      `json:"transaction_id"`
	ProductId     int64      `json:"product_id"`
	Quantity      int32      `json:"quantity"`
	Price         int64      `json:"price"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at" gorm:"index"`
}
