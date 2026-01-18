PACKAGE = github.com/wonjinsin/simple-chatbot
CUSTOM_OS = ${GOOS}
BASE_PATH = $(shell pwd)
BIN = $(BASE_PATH)/bin
BINARY_NAME = bin/server
MAIN = $(BASE_PATH)/cmd/server/main.go
GOLINT = $(BIN)/golint
MOCK = $(BIN)/mockgen
PKG_LIST = $(shell cd $(BASE_PATH) && cat pkg.list)


ifneq (, $(CUSTOM_OS))
	OS ?= $(CUSTOM_OS)
else
	OS ?= $(shell uname | awk '{print tolower($0)}')
endif

tool:
	@echo "Installing tools from go.mod..."
	@GOBIN=$(BIN) go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint
	@GOBIN=$(BIN) go install go.uber.org/mock/mockgen
	@GOBIN=$(BIN) go install github.com/golangci/golines@latest
	@echo "Tools installed successfully!"

.PHONY: build
build:
	GOOS=$(OS) go build -o $(BINARY_NAME) $(MAIN)

.PHONY: vet
vet:
	go vet

.PHONY: fmt
fmt:
	go fmt

.PHONY: lint
lint:
	$(BIN)/golangci-lint run

.PHONY: test
test: build-mocks
	go test -v -cover ./...

.PHONY: test-all
test-all: test vet fmt lint

.PHONY: build-mocks
build-mocks:
	$(MOCK) -source=internal/usecase/service.go -destination=mock/mock_service.go -package=mock
	$(MOCK) -source=internal/repository/repository.go -destination=mock/mock_repository.go -package=mock

.PHONY: init
init: 
	go mod init $(PACKAGE)

.PHONY: tidy
tidy: 
	go mod tidy

.PHONY: vendor
vendor:
	go mod vendor

# Infrastructure commands
.PHONY: infra-up
infra-up:
	docker compose up -d

.PHONY: infra-down
infra-down:
	docker compose down

# Migration command
.PHONY: migrate-up
migrate-up:
	go run cmd/migrate/main.go up

.PHONY: migrate-down
migrate-down:
	go run cmd/migrate/main.go down

.PHONY: migrate-version
migrate-version:
	go run cmd/migrate/main.go version

.PHONY: start
start: build 
	@$(BINARY_NAME)

.PHONY: all
all: tool init tidy vendor build

.PHONY: clean
clean:; $(info cleaning…) @ 
	@rm -rf vendor mock bin
	@rm -rf go.mod go.sum pkg.list
