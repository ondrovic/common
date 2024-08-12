export GO111MODULE=on
GOOS := $(shell go env GOOS)
VERSION := $(shell git describe --tags --always)
BUILD_FLAGS := -ldflags="-X 'main.version=$(VERSION)'"

# determins the variables based on GOOS 
ifeq ($(GOOS), windows)
    RM = del /Q
	HOME = $(shell echo %USERPROFILE%)
	CONFIG_PATH = $(subst  ,,$(HOME)\.golangci.yaml)
else
    RM = rm -f
	HOME = $(shell echo $$HOME)
	CONFIG_PATH = $(HOME)/.golangci.yaml
endif

check-quality:
	@make tidy
	@make fmt
	@make vet
	@make lint

lint:
	golangci-lint run --config="$(CONFIG_PATH)" ./...

vet:
	go vet ./...

fmt:
	go fmt ./...

tidy:
	go mod tidy

update-packages:
	go get -u all

update-common:
	go get github.com/ondrovic/common@latest

test:
	make tidy
	go test -v -timeout 10m ./... -coverprofile=coverage.out -json > report.json || (echo "Tests failed. See report.json for details." && exit 1)

coverage:
	make test
	go tool cover -html=coverage.out -o coverage.html

build: 
	go build $(BUILD_FLAGS)

all:
	make check-quality
	make test
	make build

clean:
	go clean
	$(RM) coverage*
	$(RM) report*
	$(RM) lint*

vendor:
	go mod vendor

