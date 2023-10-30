package main

import (
	"fmt"
	"log"
	"net"
	"os"

	productSvc "github.com/dedihartono801/transaction-svc/cmd/grpc/client/product"
	"github.com/dedihartono801/transaction-svc/database"
	"github.com/dedihartono801/transaction-svc/internal/app/repository"
	"github.com/dedihartono801/transaction-svc/internal/app/usecase/transaction"
	grpcHandler "github.com/dedihartono801/transaction-svc/internal/delivery/grpc"
	"github.com/dedihartono801/transaction-svc/pkg/identifier"
	pb "github.com/dedihartono801/transaction-svc/pkg/protobuf"
	"github.com/dedihartono801/transaction-svc/pkg/validator"
	validatorv10 "github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	mysql := database.InitMysql()

	productSvcClient := productSvc.InitServiceClient()
	identifier := identifier.NewIdentifier()
	validator := validator.NewValidator(validatorv10.New())
	dbTransactionRepository := repository.NewDbTransactionRepository(mysql)
	transactionRepository := repository.NewTransactionRepository(mysql)
	transactionService := transaction.NewGrpcTransactionService(transactionRepository, dbTransactionRepository, productSvcClient, validator, identifier)
	transactionHandler := grpcHandler.TransactionHandler{Service: transactionService}

	lis, err := net.Listen("tcp", os.Getenv("GRPC_PORT"))
	if err != nil {
		log.Fatalln("Failed to listing:", err)
	}

	fmt.Println("GRPC Svc on", os.Getenv("GRPC_PORT"))

	opts := []grpc.ServerOption{}
	tls := false //in this example we do not use ssl

	// if using ssl server side
	if tls {
		certFile := "" //example ssl path
		kefFile := ""  //example ssl path

		creds, err := credentials.NewServerTLSFromFile(certFile, kefFile)

		if err != nil {
			log.Fatalf("Failed loading certificates: %v\n", err)
		}

		opts = append(opts, grpc.Creds(creds))
	}

	grpcServer := grpc.NewServer(opts...)

	pb.RegisterTransactionServiceServer(grpcServer, &transactionHandler)
	reflection.Register(grpcServer)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln("Failed to serve:", err)
	}
}
