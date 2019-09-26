.PHONY: print gen_proto vet_server run_server

# Golang Flags
GOPATH ?= $(shell go env GOPATH)
GO=go

Server=cmd/server/server.go
Client=cmd/client/client.go

print:
	@echo print gen_proto vet_server run_server

vet_server: # run go vet
	@echo Run go vet
	$(GO) vet $(Server)

vet_client: # run go vet
	@echo Run go vet
	$(GO) vet $(Client)

run_server: vet_server
	@echo Run server
	$(GO) run $(Server)

run_client: vet_client
	@echo Run client
	$(GO) run $(Client)

build_server: vet_server
	@echo Build Server
	$(GO) build -tags=jsoniter .

gen_proto:
	@echo generator protobuf
	protoc --go_out=plugins=grpc:pkg/rpc_service -I protocol message.proto
