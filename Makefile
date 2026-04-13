.PHONY: build run test test-race test-integration test-coverage test-bench lint fmt generate clean docker/build docker/run precommit/install precommit/run

APP_NAME := sample-redis-app
CMD_DIR := ./cmd/api
COVERAGE_MIN := 95

build:
	go build -o bin/$(APP_NAME) $(CMD_DIR)

run:
	go run $(CMD_DIR)

test:
	go test ./...

test-race:
	go test -race ./...

test-integration:
	go test -tags=integration -race ./...

test-coverage:
	go test -race -coverprofile=coverage.out $$(go list ./... | grep -v /cmd/)
	go tool cover -func=coverage.out
	@echo "--- Checking coverage >= $(COVERAGE_MIN)% ---"
	@total=$$(go tool cover -func=coverage.out | grep ^total: | awk '{print $$3}' | tr -d '%'); \
	if [ $$(echo "$$total < $(COVERAGE_MIN)" | bc) -eq 1 ]; then \
		echo "FAIL: coverage $$total% < $(COVERAGE_MIN)%"; exit 1; \
	else \
		echo "OK: coverage $$total%"; \
	fi

test-bench:
	go test -bench=. -benchmem ./...

lint:
	go vet ./...
	@which staticcheck > /dev/null 2>&1 && staticcheck ./... || echo "staticcheck not installed, skipping"

fmt:
	gofmt -w .
	@which goimports > /dev/null 2>&1 && goimports -w . || echo "goimports not installed, skipping"

generate:
	go generate ./...

clean:
	rm -rf bin/ coverage.out

docker/build:
	docker build --tag $(APP_NAME):latest .

docker/run:
	docker run -p 8080:8080 $(APP_NAME):latest

precommit/install:
	pre-commit install

precommit/run:
	pre-commit run --all-files
