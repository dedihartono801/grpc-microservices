ifeq ($(OS), Windows_NT)
	PACKAGE = $(shell (Get-Content go.mod -head 1).Split(" ")[1])
else
	PACKAGE = $(shell head -1 go.mod | awk '{print $$2}')
endif

migration-up:
	migrate -database "mysql://root:root@tcp(localhost:3306)/product_db" -path migrations up

migration-down:
	migrate -database "mysql://root:root@tcp(localhost:3306)/product_db" -path migrations down

migration $$(enter):
	@read -p "Migration name:" migration_name; \
	migrate create -ext sql -dir migrations $$migration_name

proto:
	protoc -I ./pkg/protobuf --go_opt=module=${PACKAGE} --go_out=. --go-grpc_opt=module=${PACKAGE} --go-grpc_opt=require_unimplemented_servers=false --go-grpc_out=. ./pkg/protobuf/*.proto

run:
	go run cmd/grpc/main.go

format:
	go fmt ./...

run-product-svc:
	DC_APP_ENV=dev $(GOPATH)/bin/reflex -s -r '\.go$$' make format run

test-cov:
	go test -coverprofile=cover.out ./... && go tool cover -html=cover.out -o cover.html