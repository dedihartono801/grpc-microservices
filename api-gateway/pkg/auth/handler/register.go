package handler

import (
	"context"
	"net/http"

	"github.com/dedihartono801/api-gateway/pkg/auth/pb"
	"github.com/dedihartono801/api-gateway/pkg/customstatus"
	"github.com/dedihartono801/api-gateway/pkg/helper"
	"github.com/gin-gonic/gin"
)

type RegisterRequestBody struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(ctx *gin.Context, c pb.UserServiceClient) {
	b := RegisterRequestBody{}

	if err := ctx.BindJSON(&b); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	res, err := c.Register(context.Background(), &pb.RegisterRequest{
		Username: b.Username,
		Email:    b.Email,
		Password: b.Password,
	})

	if err != nil {
		response := helper.CustomResponse(ctx, customstatus.ErrInternalServerFailed.Message, customstatus.ErrInternalServerError.Message, nil)
		ctx.JSON(http.StatusBadGateway, response)
		return
	}

	if res.Status != http.StatusCreated {
		response := helper.CustomResponse(ctx, customstatus.ErrInternalServerFailed.Message, res.Error, nil)
		ctx.JSON(int(res.Status), response)
		return
	}
	response := helper.CustomResponse(ctx, customstatus.StatusCreated.Message, "", res.Data)
	ctx.JSON(int(res.Status), response)
}
