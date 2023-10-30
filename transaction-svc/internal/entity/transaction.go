package entity

import (
	"time"
)

type Transaction struct {
	ID            int64      `gorm:"primaryKey;autoIncrement:true" json:"id"`
	UserId        int64      `json:"user_id"`
	TotalAmount   int64      `json:"total_amount"`
	TotalQuantity int64      `json:"total_quantity"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at" gorm:"index"`
}
