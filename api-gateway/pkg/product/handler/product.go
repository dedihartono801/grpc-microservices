package handler

import (
	"context"
	"net/http"

	"github.com/dedihartono801/api-gateway/pkg/customstatus"
	"github.com/dedihartono801/api-gateway/pkg/helper"
	"github.com/dedihartono801/api-gateway/pkg/product/pb"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ProductRequestRequestBody struct {
	ProductId int64 `json:"product_id" validate:"required"`
	Quantity  int32 `json:"quantity" validate:"required"`
}

func GetProduct(ctx *gin.Context, c pb.ProductServiceClient, redis *redis.Client) {
	res, err := c.ListProduct(context.Background(), &emptypb.Empty{})

	if err != nil {
		response := helper.CustomResponse(ctx, customstatus.ErrInternalServerFailed.Message, res.Error, nil)
		ctx.JSON(http.StatusBadGateway, response)
		return
	}

	if res.Status != http.StatusOK {
		response := helper.CustomResponse(ctx, customstatus.ErrInternalServerFailed.Message, res.Error, nil)
		ctx.JSON(int(res.Status), response)
		return
	}

	// helper.RedisSet(ctx, "products", res.Data, *redis)

	response := helper.CustomResponse(ctx, customstatus.StatusOk.Message, "", res.Data)
	ctx.JSON(int(res.Status), response)
}
