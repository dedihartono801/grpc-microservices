package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dedihartono801/api-gateway/pkg/customstatus"
	"github.com/dedihartono801/api-gateway/pkg/helper"
	"github.com/dedihartono801/api-gateway/pkg/transaction/pb"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type TransactionItems struct {
	Items []Items `json:"items" validate:"required"`
}

type Items struct {
	ProductId int64 `json:"product_id" validate:"required"`
	Quantity  int32 `json:"quantity" validate:"required"`
}

func CreateTransaction(ctx *gin.Context, c pb.TransactionServiceClient) {
	b := TransactionItems{}
	user_id := ctx.GetInt64("user_id")

	if err := ctx.BindJSON(&b); err != nil {
		response := helper.CustomResponse(ctx, customstatus.ErrInternalServerFailed.Message, err.Error(), nil)
		ctx.JSON(http.StatusBadGateway, response)
		return
	}

	dt := make([]*pb.Items, len(b.Items))
	for i, val := range b.Items {
		dt[i] = &pb.Items{
			ProductId: val.ProductId,
			Quantity:  val.Quantity,
		}
	}

	res, err := c.Transaction(context.Background(), &pb.TransactionRequest{
		UserId: user_id,
		Items:  dt,
	})

	if err != nil {
		response := helper.CustomResponse(ctx, customstatus.ErrInternalServerFailed.Message, err.Error(), nil)
		ctx.JSON(http.StatusBadGateway, response)
		return
	}

	if res.Error != "" {
		response := helper.CustomResponse(ctx, customstatus.ErrInternalServerFailed.Message, res.Error, nil)
		ctx.JSON(http.StatusBadGateway, response)
		return
	}

	if res.Status != http.StatusCreated {
		response := helper.CustomResponse(ctx, customstatus.ErrInternalServerFailed.Message, res.Error, nil)
		ctx.JSON(int(res.Status), response)
		return
	}
	response := helper.CustomResponse(ctx, customstatus.StatusOk.Message, "", res.TransactionId)
	ctx.JSON(int(res.Status), response)
}

func GetDetailTransaction(ctx *gin.Context, c pb.TransactionServiceClient, redis *redis.Client) {
	trxId := ctx.Param("transactionId")
	trxIdInt, _ := strconv.Atoi(trxId)
	transaction, err := helper.RedisGetTrx(ctx, "transaction-"+trxId, *redis)
	if err != nil {
		response := helper.CustomResponse(ctx, customstatus.ErrInternalServerFailed.Message, customstatus.ErrInternalServerError.Message, nil)
		ctx.JSON(http.StatusBadGateway, response)
		return
	}

	if transaction != nil && err == nil {
		fmt.Println("Get from redis")
		response := helper.CustomResponse(ctx, customstatus.StatusOk.Message, "", transaction)
		ctx.JSON(http.StatusOK, response)
		return
	}

	fmt.Println("Get from non redis")
	res, err := c.GetDetailTransaction(context.Background(), &pb.GetDetailTransactionRequest{
		TransactionId: int64(trxIdInt),
	})

	if err != nil {
		response := helper.CustomResponse(ctx, customstatus.ErrInternalServerFailed.Message, customstatus.ErrNotFound.Message, nil)
		ctx.JSON(http.StatusNotFound, response)
		return
	}

	helper.RedisSetTrx(ctx, "transaction-"+trxId, res.Data, *redis)

	if res.Status != http.StatusOK {
		response := helper.CustomResponse(ctx, customstatus.ErrInternalServerFailed.Message, res.Error, nil)
		ctx.JSON(int(res.Status), response)
		return
	}
	response := helper.CustomResponse(ctx, customstatus.StatusCreated.Message, "", res.Data)
	ctx.JSON(int(res.Status), response)
}
