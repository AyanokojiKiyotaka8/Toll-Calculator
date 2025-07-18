obu:
	@go build -o bin/obu ./obu
	@./bin/obu

rec:
	@go build -o bin/receiver ./data_receiver
	@./bin/receiver

calc:
	@go build -o bin/calculator ./distance_calculator
	@./bin/calculator

agg:
	@go build -o bin/aggregator ./aggregator
	@./bin/aggregator

PROTOC_GEN_GO := $(shell go env GOPATH)/bin/protoc-gen-go
PROTOC_GEN_GRPC := $(shell go env GOPATH)/bin/protoc-gen-go-grpc

proto:
	@protoc --plugin=protoc-gen-go=$(PROTOC_GEN_GO) \
	--plugin=protoc-gen-go-grpc=$(PROTOC_GEN_GRPC) \
	--go_out=. --go_opt=paths=source_relative \
	--go-grpc_out=. --go-grpc_opt=paths=source_relative \
	types/ptypes.proto

.PHONY: obu rec calc agg proto