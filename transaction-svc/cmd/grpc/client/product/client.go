package user

import (
	"fmt"
	"log"

	productPb "github.com/dedihartono801/transaction-svc/cmd/grpc/client/product/pb"
	"github.com/dedihartono801/transaction-svc/pkg/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type ServiceClient struct {
	Product productPb.ProductServiceClient
}

func InitServiceClient() ServiceClient {
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

	fmt.Println("Connect to product service")

	return ServiceClient{
		Product: productPb.NewProductServiceClient(cc),
	}
}
