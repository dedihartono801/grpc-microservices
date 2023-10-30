package entity

import (
	"time"
)

type Product struct {
	ID        int64      `gorm:"primaryKey;autoIncrement:true" json:"id"`
	Name      string     `json:"name" validate:"required"`
	Stock     int32      `json:"stock" validate:"required"`
	Price     int64      `json:"price" validate:"required"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"index"`
}
