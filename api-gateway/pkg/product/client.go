package product

import (
	"fmt"
	"log"

	"github.com/dedihartono801/api-gateway/pkg/config"
	pb "github.com/dedihartono801/api-gateway/pkg/product/pb"
	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type ServiceClient struct {
	Product pb.ProductServiceClient
	Redis   *redis.Client
}

func InitServiceClient(redis *redis.Client) ServiceClient {
	opts := []grpc.DialOption{}
	tls := false // using WithInsecure() because no SSL running

	if tls {
		certFile := ""

		creds, err := credentials.NewClientTLSFromFile(certFile, "")

		if err != nil {
			log.Fatalf("Error while loading CA trust certificates: %v\n", err)
		}

		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	// using WithSecure() because no SSL running
	cc, err := grpc.Dial(config.GetEnv("PRODUCT_SVC_ADDR"), opts...)

	if err != nil {
		fmt.Println("Could not connect:", err)
	}

	return ServiceClient{
		Product: pb.NewProductServiceClient(cc),
		Redis:   redis,
	}
}
