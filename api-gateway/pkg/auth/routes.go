package auth

import (
	"github.com/dedihartono801/api-gateway/pkg/auth/handler"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) ServiceClient {
	svc := InitServiceClient()

	routes := r.Group("/user")
	routes.POST("/register", svc.Register)
	routes.POST("/login", svc.Login)

	return svc
}

func (svc *ServiceClient) Register(ctx *gin.Context) {
	handler.Register(ctx, svc.User)
}

func (svc *ServiceClient) Login(ctx *gin.Context) {
	handler.Login(ctx, svc.User)
}
