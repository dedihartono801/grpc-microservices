package transaction

import (
	"github.com/dedihartono801/api-gateway/pkg/auth"
	"github.com/dedihartono801/api-gateway/pkg/transaction/handler"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func RegisterRoutes(r *gin.Engine, authSvc *auth.ServiceClient, redis *redis.Client) ServiceClient {
	a := auth.InitAuthMiddleware(authSvc)
	svc := InitServiceClient(redis)

	routes := r.Group("/transaction")
	routes.Use(a.AuthRequired)
	routes.POST("", svc.CreateTransaction)
	routes.GET("/:transactionId", svc.GetDetailTransaction)

	return svc
}

func (svc *ServiceClient) CreateTransaction(ctx *gin.Context) {
	handler.CreateTransaction(ctx, svc.Transaction)
}

func (svc *ServiceClient) GetDetailTransaction(ctx *gin.Context) {
	handler.GetDetailTransaction(ctx, svc.Transaction, svc.Redis)
}
