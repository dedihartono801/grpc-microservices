package main

import (
	"log"

	"github.com/dedihartono801/api-gateway/pkg/auth"
	"github.com/dedihartono801/api-gateway/pkg/config"
	"github.com/dedihartono801/api-gateway/pkg/logger"
	"github.com/dedihartono801/api-gateway/pkg/product"
	"github.com/dedihartono801/api-gateway/pkg/transaction"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	if err != nil {
		log.Fatalln("Failed at config", err)
	}

	redis := redis.NewClient(&redis.Options{
		Addr:     config.GetEnv("REDIS_HOST"), // Redis server address
		Password: "",                          // No password by default
		DB:       0,                           // Default DB
	})

	logger.InitLogger()

	r := gin.Default()

	logger.CreateLog(r)

	auth := auth.RegisterRoutes(r)
	product.RegisterRoutes(r, redis)
	transaction.RegisterRoutes(r, &auth, redis)

	r.Run(config.GetEnv("PORT"))
}
