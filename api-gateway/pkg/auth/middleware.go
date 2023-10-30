package auth

import (
	"context"
	"net/http"

	pb "github.com/dedihartono801/api-gateway/pkg/auth/pb"
	"github.com/dedihartono801/api-gateway/pkg/customstatus"
	"github.com/dedihartono801/api-gateway/pkg/helper"
	"github.com/gin-gonic/gin"
)

type AuthMiddlewareConfig struct {
	svc *ServiceClient
}

func InitAuthMiddleware(svc *ServiceClient) AuthMiddlewareConfig {
	return AuthMiddlewareConfig{svc}
}

func (c *AuthMiddlewareConfig) AuthRequired(ctx *gin.Context) {
	authorization := ctx.Request.Header.Get("authorization")

	res, err := c.svc.User.Validate(context.Background(), &pb.ValidateRequest{
		Token: authorization,
	})

	if err != nil || res.Status != http.StatusOK {
		response := helper.CustomResponse(ctx, customstatus.ErrInternalServerFailed.Message, res.Error, nil)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}

	ctx.Set("user_id", int64(res.UserId))

	ctx.Next()
}
