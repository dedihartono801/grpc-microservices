package product

import (
	"testing"

	repoMock "github.com/dedihartono801/product-svc/internal/app/repository"
	"github.com/dedihartono801/product-svc/internal/entity"
	"github.com/dedihartono801/product-svc/pkg/identifier"
	"github.com/dedihartono801/product-svc/pkg/validator"
	validatorv10 "github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestListProduct(t *testing.T) {
	productRepo := repoMock.NewMockProductRepository()
	identifier := identifier.NewIdentifier()
	validator := validator.NewValidator(validatorv10.New())
	srv := NewGrpcProductService(productRepo, validator, identifier)

	expected := []*entity.Product{
		{
			ID:    1,
			Name:  "Kopi ABC dudu",
			Stock: 100,
			Price: 2000,
		},
	}

	// Define test cases
	testCases := []struct {
		name       string
		expected   []*entity.Product
		wantErr    bool
		statusCode int
	}{
		{
			name:       "Success List Product",
			expected:   expected,
			wantErr:    false,
			statusCode: 201,
		},
	}

	// Run tests
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := srv.ListProduct()
			assert.Equal(t, tc.expected, actual, "Expected and actual data should be equal")
			assert.Equal(t, tc.wantErr, err != nil, "Expected error and actual error should be equal")
		})
	}

}
