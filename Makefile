COVERAGE_FILE ?= coverage.out

.PHONY: build
build: build_bot build_api

.PHONY: build_bot
build_bot:
	@echo "Выполняется go build для таргета bot"
	@mkdir -p bin
	@go build -o ./bin/bot ./cmd/doto-bot/

.PHONY: build_api
build_api:
	@echo "Выполняется go build для таргета api"
	@mkdir -p bin
	@go build -o ./bin/api ./cmd/doto/

## test: run all tests
.PHONY: test
test:
	@go test -coverpkg='github.com/es-debug/backend-academy-2024-go-template/...' --race -count=1 -coverprofile='$(COVERAGE_FILE)' ./...
	@go tool cover -func='$(COVERAGE_FILE)' | grep ^total | tr -s '\t'
	@go tool cover -html='$(COVERAGE_FILE)' -o coverage.html && xdg-open coverage.html

.PHONY: lint
lint:
	@if ! command -v 'golangci-lint' &> /dev/null; then \
			echo "Please install golangci-lint!"; exit 1; \
		fi;
	@golangci-lint -v run --fix ./...

.PHONY: generate
generate:
	@if ! command -v 'oapi-codegen' &> /dev/null; then \
		echo "Please install oapi-codegen!"; exit 1; \
	fi;
	@mkdir -p internal/clients/doto
	@mkdir -p api/
	@oapi-codegen --config oapi.doto-client.yaml oapi/doto-api.yaml
	@oapi-codegen --config oapi.doto-server.yaml oapi/doto-api.yaml

.PHONY: clean
clean:
	@rm -rf bin
