# Makefile — fallback for environments without Task (go-task) installed.
# Prefer using Taskfile.yml via `task <target>` when possible.

GOOS        := linux
GOARCH      := arm64
CGO_ENABLED := 0
BUILD_FLAGS := -tags lambda.norpc -trimpath -ldflags="-s -w"
BIN_DIR     := bin

.PHONY: build build-func test-unit test-integration test-all lint audit \
        package deploy-dev deploy-staging deploy-prod \
        local-up local-down clean mod-tidy

build:
	@for dir in cmd/*/; do \
		func=$$(basename "$$dir"); \
		echo "Building $$func..."; \
		GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=$(CGO_ENABLED) \
			go build $(BUILD_FLAGS) -o $(BIN_DIR)/$$func/bootstrap ./cmd/$$func/; \
	done

build-func:
	@test -n "$(FUNC)" || (echo "Usage: make build-func FUNC=myhandler" && exit 1)
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=$(CGO_ENABLED) \
		go build $(BUILD_FLAGS) -o $(BIN_DIR)/$(FUNC)/bootstrap ./cmd/$(FUNC)/

test-unit:
	go test -short -count=1 -race ./...

test-integration:
	go test -run Integration -count=1 -race ./...

test-all:
	go test -count=1 -race ./...

lint:
	golangci-lint run ./...

audit:
	govulncheck ./...

package: build
	@for dir in $(BIN_DIR)/*/; do \
		func=$$(basename "$$dir"); \
		echo "Packaging $$func..."; \
		(cd "$$dir" && zip "../$$func.zip" bootstrap); \
	done

deploy-dev: package
	bash scripts/deploy.sh dev

deploy-staging: package
	bash scripts/deploy.sh staging

deploy-prod: package
	bash scripts/deploy.sh prod

local-up:
	docker-compose up -d

local-down:
	docker-compose down

clean:
	rm -rf $(BIN_DIR)
	go clean -cache -testcache

mod-tidy:
	go mod tidy
	go mod verify
