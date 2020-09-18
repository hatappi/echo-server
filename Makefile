GO=go
DOCKER=docker

DOCKER_TAG ?= v0.4
DOCKER_IMAGE_NAME=hatappi/echo-server:${DOCKER_TAG}

.PHONY: build
build:
	@${GO} build -o dist/echo-server main.go

.PHONY: build-image
build-image:
	${DOCKER} build -t ${DOCKER_IMAGE_NAME} .

.PHONY: push-image
push-build: build-image
	${DOCKER} push ${DOCKER_IMAGE_NAME}
