GO=go
DOCKER=docker

DOCKER_TAG ?= v0.9.8
DOCKER_IMAGE_NAME=hatappi/echo-server:${DOCKER_TAG}

.PHONY: build
build:
	@${GO} build -o dist/echo-server main.go

.PHONY: build-image
build-image:
	${DOCKER} build -t ${DOCKER_IMAGE_NAME} .

.PHONY: push-image
push-image: build-image
	${DOCKER} push ${DOCKER_IMAGE_NAME}

.PHONY: protoc
protoc:
	protoc --go_out=. --go_opt=paths=source_relative \
	    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	    pb/echo.proto

