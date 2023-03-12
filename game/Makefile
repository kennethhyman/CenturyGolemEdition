SHELL := /bin/bash

.SILENT: build_protos build_game test_game


all: build_protos build_game

build_game:
	echo "building game..."
	go build .

build_protos:
	echo "building grpc objects..."
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		grpc/CenturyGolemEdition.proto

test_game: build_game
	echo "running tests..."
	go test ./models/ -v
