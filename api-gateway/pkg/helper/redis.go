package helper

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	redisClient "github.com/go-redis/redis/v8"
)

type Item struct {
	ID        int `json:"id"`
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
	Price     int `json:"price"`
}

type Transaction struct {
	ID            int    `json:"id"`
	UserID        int    `json:"user_id"`
	TotalAmount   int    `json:"total_amount"`
	TotalQuantity int    `json:"total_quantity"`
	Items         []Item `json:"items"`
}

func RedisGetTrx(ctx *gin.Context, key string, redis redisClient.Client) (*Transaction, error) {
	value, err := redis.Get(context.Background(), key).Result()
	if err != nil {
		if err != redisClient.Nil {
			return nil, err
		}
	} else {
		var transaction Transaction
		if err := json.Unmarshal([]byte(value), &transaction); err != nil {
			fmt.Println("Error:", err)
		}
		return &transaction, nil
	}
	return nil, nil
}

func RedisSetTrx(ctx *gin.Context, key string, data interface{}, redis redisClient.Client) error {
	data, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	expiration := 30 * time.Minute
	err = redis.Set(context.Background(), key, data, expiration).Err()
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	return nil
}
