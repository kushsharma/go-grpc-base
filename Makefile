.ONESHELL:
.SHELLFLAGS := -eu -o pipefail -c
.DELETE_ON_ERROR:
MAKEFLAGS += --no-builtin-rules
VERSION=`cat version`
BUILD=`date +%FT%T%z`
#COMMIT=`git rev-parse HEAD`
COMMIT=`date +%FT%T%z`
EXECUTABLE="gbase"

# first command used as the default one if only `make` is used
all: build

install:
	@echo "> installing dependencies"
	go get -u github.com/golang/protobuf/{proto,protoc-gen-go}
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc
	go get -u google.golang.org/grpc
	go get -u github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
	go get -u github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
	go get -u github.com/bufbuild/buf/cmd/buf

.PHONY: build test clean generate dist init build_linux build_mac install generate

generate:
	@echo "> building assets"
	@buf generate --path ./proto

# app.pb.go: ./proto/api/v1/app.proto
# 	@echo "> building protos"
# 	@protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./proto/api/v1/app.proto

# build: app.pb.go
build: generate
	@echo "> building binary"
	@go build -ldflags "-X 'main.Version=${VERSION}' -X 'main.Build=${BUILD}' -X 'main.AppName=${EXECUTABLE}'" -o ${EXECUTABLE} ./main.go
	@echo "> build complete"

run: build
	@./${EXECUTABLE}

clean:
	@rm -rf ${EXECUTABLE} dist/

build_nix:
	@env GOOS=linux GOARCH=amd64 go build -ldflags "-X 'main.Version=${VERSION}' -X 'main.Build=${BUILD}' -X 'main.AppName=${EXECUTABLE}'" -o ${EXECUTABLE} ./main.go

build_mac:
	@env GOOS=darwin GOARCH=amd64 go build -ldflags "-X 'main.Version=${VERSION}' -X 'main.Build=${BUILD}' -X 'main.AppName=${EXECUTABLE}'" -o ${EXECUTABLE} ./main.go
