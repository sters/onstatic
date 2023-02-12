
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
	$(GOBIN)/golangci-lint run -v ./...

.PHONY: lint-fix
lint-fix:
	$(GOBIN)/golangci-lint run --fix -v ./...

.PHONY: test
test:
	go test -v -race ./...

.PHONY: cover
cover:
	go test -v -race -coverpkg=./... -coverprofile=coverage.txt ./...

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: build-proto
build-proto:
	protoc \
		--go_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_out=. \
		--go-grpc_opt=paths=source_relative \
		pluginapi/plugin.proto

.PHONY: build-example-plugin
build-example-plugin:
	@for dir in $(shell ls plugins); do \
		echo "Build: plugins/$${dir}"; \
		go build -o plugins/$${dir}/$${dir} plugins/$${dir}/main.go; \
	done
