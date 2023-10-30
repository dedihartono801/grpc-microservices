package repository

import (
	"github.com/dedihartono801/product-svc/internal/entity"
	"gorm.io/gorm"
)

type ProductRepository interface {
	ListProduct() ([]*entity.Product, error)
	GetProductById(id int64) (*entity.Product, error)
	UpdateStockProduct(product *entity.Product) error
}

type productRepository struct {
	database *gorm.DB
}

func NewProductRepository(database *gorm.DB) ProductRepository {
	return &productRepository{database}
}

func (r *productRepository) ListProduct() ([]*entity.Product, error) {
	products := []*entity.Product{}
	err := r.database.Table("products").Scan(&products).Error
	return products, err
}

func (r *productRepository) GetProductById(id int64) (*entity.Product, error) {
	product := entity.Product{}
	err := r.database.Table("products").Where("id = ?", id).Find(&product).Error
	return &product, err
}

func (r *productRepository) UpdateStockProduct(product *entity.Product) error {

	return r.database.Table("products").Save(product).Error
}
