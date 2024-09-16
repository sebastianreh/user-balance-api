LINT_VERSION = v1.61.0
.PHONY: code-format-checks
code-format-check:
	@unformatted_files="$$(gofmt -l .)" \
	&& test -z "$$unformatted_files" || ( printf "Unformatted files: \n$${unformatted_files}\nRun make code-format\n"; exit 1 )

lint:
	golangci-lint run --config golangci.yml

lint-install:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin $(LINT_VERSION)
.PHONY: code-format

code-format:
	goimports -l -w .
	gofmt -l -w .

run-test:
	go test ./...

create-swag-docs:
	swag init -g="./cmd/main/main.go" --dir="./,./internal/interfaces/http" -o="./docs/swagger"

build-server-image:
	DOCKER_BUILDKIT=1  docker build --force-rm -t user-balance-api --no-cache .

start-compose:
	docker-compose up -d

down-compose:
	docker-compose down

create-migration-csv:
	python3 scripts/generate_transactions/generate_users_and_transactions.py