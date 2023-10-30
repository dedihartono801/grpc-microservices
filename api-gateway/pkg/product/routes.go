package product

import (
	"github.com/dedihartono801/api-gateway/pkg/product/handler"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func RegisterRoutes(r *gin.Engine, redis *redis.Client) ServiceClient {
	svc := InitServiceClient(redis)

	r.GET("/products", svc.GetProducts)

	return svc
}

func (svc *ServiceClient) GetProducts(ctx *gin.Context) {
	handler.GetProduct(ctx, svc.Product, svc.Redis)
}
