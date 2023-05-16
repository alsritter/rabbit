GOHOSTOS:=$(shell go env GOHOSTOS)
GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)

ifeq ($(GOHOSTOS), windows)
	#the `find.exe` is different from `find` in bash/shell.
	#to see https://docs.microsoft.com/en-us/windows-server/administration/windows-commands/find.
	#changed to use git-bash.exe to run find cli or other cli friendly, caused of every developer has a Git.
	Git_Bash= $(subst cmd\,bin\bash.exe,$(dir $(shell where git)))
	INTERNAL_PROTO_FILES=$(shell find internal -name *.proto)
	API_PROTO_FILES=$(shell find api -name *.proto)
else
	INTERNAL_PROTO_FILES=$(shell find internal -name *.proto)
	API_PROTO_FILES=$(shell find api -name *.proto)
endif

.PHONY: init
# init env
init:
	# apt install -y protobuf-compiler
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.30.0
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-errors/v2@v2.0.0-20230515030202-6d741828c2d4
	go install github.com/go-kratos/kratos/cmd/kratos/v2@v2.0.0-20230515030202-6d741828c2d4
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@v2.0.0-20230515030202-6d741828c2d4
	go install github.com/google/gnostic/cmd/protoc-gen-openapi@v0.6.8
	go install github.com/envoyproxy/protoc-gen-validate@v1.0.1
	go install github.com/google/wire/cmd/wire@v0.5.0


.PHONY: clean
# clean cache and ignore file
clean:
	rm -rf ./api/**/*.go


.PHONY: config
# generate internal proto
config:
	protoc --proto_path=./internal \
	       --proto_path=./third_party \
 	       --go_out=paths=source_relative:./internal \
	       $(INTERNAL_PROTO_FILES)

.PHONY: api
# generate api proto
api:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	protoc --proto_path=./api \
	       --proto_path=./third_party \
 	       --go_out=paths=source_relative:./api \
 	       --go-http_out=paths=source_relative:./api \
 	       --go-grpc_out=paths=source_relative:./api \
	       --openapi_out=fq_schema_naming=true,default_response=false:. \
	       $(API_PROTO_FILES)


.PHONY: validate
# generate validate proto
validate:
	protoc --proto_path=./api \
           --proto_path=./third_party \
           --go_out=paths=source_relative:./api \
           --validate_out=paths=source_relative,lang=go:./api \
           $(API_PROTO_FILES)


.PHONY: build
# build
build:
	mkdir -p bin/ && CGO_ENABLED=0 GOOS=linux go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./...

.PHONY: generate
# generate
generate:
	# go mod tidy
	# go get github.com/google/wire/cmd/wire@latest
	# go generate ./...

.PHONY: wire
# generate
wire:
	go mod tidy
	cd ./cmd/server && wire

.PHONY: err
err:
	protoc --proto_path=. \
			--proto_path=./third_party \
			--go_out=paths=source_relative:. \
			--go-errors_out=paths=source_relative:. \
			api/errcode/error.proto

.PHONY: all
# generate all
all: clean api validate err config wire

.PHONY: release
# release to server
release: all build


run:
	go run ./cmd/server/.

# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help

