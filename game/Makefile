SHELL := /bin/bash

.SILENT: build_protos build_game test_game tidy vet lint run_server run_cli clean


all: tidy vet build_protos build_game test_game

build_game:
	echo "building game..."
	mkdir -p ./bin
	go build -o ./bin ./cmd/...

build_protos:
	echo "building grpc objects..."
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		internal/core/grpc/CenturyGolemEdition.proto

test_game:
	echo "running tests..."
	go test ./... -v

tidy:
	echo "tidying go.mod..."
	go mod tidy

vet:
	echo "vetting..."
	go vet ./...

lint:
	echo "linting..."
	golangci-lint run ./...

run_server: build_game
	./bin/server

run_cli: build_game
	./bin/cli

clean:
	echo "cleaning..."
	rm -rf ./bin
