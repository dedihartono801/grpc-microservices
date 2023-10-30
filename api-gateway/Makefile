PROTO_DIR = pb

ifeq ($(OS), Windows_NT)
	PACKAGE = $(shell (Get-Content go.mod -head 1).Split(" ")[1])
else
	PACKAGE = $(shell head -1 go.mod | awk '{print $$2}')
endif

.PHONY: auth product transaction
project := auth product transaction


$(project):
	@${CHECK_DIR_CMD}
	protoc -Ipkg/$@/${PROTO_DIR} --go_opt=module=${PACKAGE} --go_out=. --go-grpc_opt=module=${PACKAGE} --go-grpc_opt=require_unimplemented_servers=false \
	--go-grpc_out=. pkg/$@/${PROTO_DIR}/*.proto

run:
	go run cmd/main.go

format:
	go fmt ./...

run-apigw-svc:
	DC_APP_ENV=dev $(GOPATH)/bin/reflex -s -r '\.go$$' make format run