package repository

import (
	"github.com/dedihartono801/product-svc/internal/entity"
)

// Define the MockAdminRepository interface
type MockProductRepository interface {
	ListProduct() ([]*entity.Product, error)
	GetProductById(id int64) (*entity.Product, error)
	UpdateStockProduct(product *entity.Product) error
}

type mockProductRepository struct{}

func NewMockProductRepository() MockProductRepository {
	return &mockProductRepository{}

}

// List returns a list of all skus in the repository
func (m *mockProductRepository) ListProduct() ([]*entity.Product, error) {
	var products []*entity.Product

	dt := &entity.Product{
		ID:    1,
		Name:  "Kopi ABC dudu",
		Stock: 100,
		Price: 2000,
	}
	products = append(products, dt)
	return products, nil
}

func (m *mockProductRepository) GetProductById(id int64) (*entity.Product, error) {

	dt := &entity.Product{
		ID:    1,
		Name:  "Kopi ABC dudu",
		Stock: 100,
		Price: 2000,
	}
	return dt, nil
}

func (m *mockProductRepository) UpdateStockProduct(product *entity.Product) error {
	return nil
}
