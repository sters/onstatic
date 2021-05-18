
export GOBIN := $(PWD)/bin
export PATH := $(GOBIN):$(PATH)

TOOLS=$(shell cat tools/tools.go | egrep '^\s_ '  | awk '{ print $$2 }')

.PHONY: bootstrap-tools
bootstrap-tools:
	@echo "Installing: " $(TOOLS)
	@go install $(TOOLS)

.PHONY: run
run:
	go run cmd/server/main.go

.PHONY: lint
lint:
	golangci-lint run -v ./...
	go-consistent -v ./...

.PHONY: lint-fix
lint-fix:
	golangci-lint run --fix -v ./...

.PHONY: test
test:
	go test -v -race ./...

.PHONY: cover
cover:
	go test -v -race -coverpkg=./... -coverprofile=coverage.txt ./...

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: build-example-plugin
build-example-plugin:
	CGO_ENABLED=1 go build -buildmode=plugin -o plugins/example/example.so plugins/example/main.go
