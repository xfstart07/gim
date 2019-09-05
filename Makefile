.PHONY: print

# Golang Flags
GOPATH ?= $(shell go env GOPATH)
GO=go

Server=cmd/server/server.go
Client=cmd/client/client.go

print:
	@echo print

vet_server: # run go vet
	@echo Run go vet
	go vet $(Server)

run_server: vet_server
	@echo Run server
	$(GO) run $(Server)
