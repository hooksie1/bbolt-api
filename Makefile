PROJECT_NAME := "bbolt-api"
PKG := "github.com/hooksie1/$(PROJECT_NAME)"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)
VERSION := $$(git describe --tags | cut -d '-' -f 1)


.PHONY: all build docker dep clean test coverage lint

all: build

lint: ## Lint the files
	@golint -set_exit_status ./...

test: ## Run unittests
	@go test ./...

coverage:
	@go test -cover ./...
	@go test ./... -coverprofile=cover.out && go tool cover -html=cover.out -o coverage.html

dep: ## Get the dependencies
	@go get -u golang.org/x/lint/golint

build: linux windows mac

linux: dep
	@CGO_ENABLED=0 GOOS=linux go build -a -ldflags "-w -X '$(PKG)/cmd.Version=$(VERSION)'" -o bbolt-api

windows: dep
	@CGO_ENABLED=0 GOOS=windows go build -a -ldflags "-w -X '$(PKG)/cmd.Version=$(VERSION)'" -o bbolt-api.exe

mac: dep
	@CGO_ENABLED=0 GOOS=darwin go build -a -ldflags "-w -X '$(PKG)/cmd.Version=$(VERSION)'" -o bbolt-api-darwin

docker: build
	@docker build -f Dockerfile.app -t hooksie1/bbolt-api:$(VERSION) .

push: docker
	@docker push hooksie1/bbolt-api:$(VERSION)

clean: ## Remove previous build
	git clean -fd
	git clean -fx
	git reset --hard

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
