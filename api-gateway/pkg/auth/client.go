package auth

import (
	"fmt"
	"log"

	pb "github.com/dedihartono801/api-gateway/pkg/auth/pb"
	"github.com/dedihartono801/api-gateway/pkg/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type ServiceClient struct {
	User pb.UserServiceClient
}

func InitServiceClient() ServiceClient {
	opts := []grpc.DialOption{}
	tls := false // using WithInsecure() because no SSL running

	if tls {
		certFile := "t"

		creds, err := credentials.NewClientTLSFromFile(certFile, "")

		if err != nil {
			log.Fatalf("Error while loading CA trust certificates: %v\n", err)
		}

		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	// using WithSecure() because no SSL running
	cc, err := grpc.Dial(config.GetEnv("AUTH_SVC_ADDR"), opts...)

	if err != nil {
		fmt.Println("Could not connect:", err)
	}

	return ServiceClient{
		User: pb.NewUserServiceClient(cc),
	}
}
