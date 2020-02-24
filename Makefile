.PHONY: dev test image

APP_NAME := urbs-setting
APP_VERSION := $(shell git describe --tags --always --match "v[0-9]*")
APP_PATH := $(shell echo ${PWD} | sed -e "s\#${GOPATH}/src/\#\#g")

dev:
	@CONFIG_FILE_PATH=${PWD}/config/default.yml go run main.go

test:
	@CONFIG_FILE_PATH=${PWD}/config/test-local.yml go test -v ./...

BUILD_TIME := $(shell date -u +"%FT%TZ")
BUILD_COMMIT := $(shell git rev-parse HEAD)

.PHONY: build build-tool
build:
	@mkdir -p ./dist
	GO111MODULE=on go build -ldflags "-X ${APP_PATH}/src/api.AppName=${APP_NAME} \
	-X ${APP_PATH}/src/api.AppVersion=${APP_VERSION} \
	-X ${APP_PATH}/src/api.BuildTime=${BUILD_TIME} \
	-X ${APP_PATH}/src/api.GitSHA1=${BUILD_COMMIT}" \
	-o ./dist/urbs-setting main.go
build-tool:
	@mkdir -p ./dist
	GO111MODULE=on go build -ldflags "-X ${APP_PATH}/src/api.AppName=${APP_NAME} \
	-X ${APP_PATH}/src/api.AppVersion=${APP_VERSION} \
	-X ${APP_PATH}/src/api.BuildTime=${BUILD_TIME} \
	-X ${APP_PATH}/src/api.GitSHA1=${BUILD_COMMIT}" \
	-o ./dist/sql-cli cmd/sql_cli/sql_cli.go

PKG_LIST := $(shell go list ./... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/)

.PHONY: lint
lint:
	@hash golint > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		go get -u golang.org/x/lint/golint; \
	fi
	@golint -set_exit_status ${PKG_LIST}

.PHONY: fmt-check
fmt-check:
	test -z "$(shell gofmt -d -e ${GO_FILES})"

.PHONY: misspell-check
misspell-check:
	@hash misspell > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		go get -u github.com/client9/misspell/cmd/misspell; \
	fi
	@misspell -error $(GO_FILES)

.PHONY: coverhtml
coverhtml:
	@mkdir -p coverage
	@CONFIG_FILE_PATH=${PWD}/config/test-local.yml go test -coverprofile=coverage/cover.out ./...
	@go tool cover -html=coverage/cover.out -o coverage/coverage.html
	@go tool cover -func=coverage/cover.out | tail -n 1

DOCKER_IMAGE_TAG := ${APP_NAME}:${APP_VERSION}
.PHONY: image
image:
	docker build --rm -t ${DOCKER_IMAGE_TAG} .
