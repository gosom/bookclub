#!make
default: help

# generate help info from comments: thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help: ## help information about make commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

test: ## runs the unit tests
	go test -v -timeout 5m ./...

test-cover: ## outputs the unittest coverage statistics
	go test -v -timeout 5m ./... -coverprofile coverage.out
	go tool cover -func coverage.out
	rm coverage.out

test-cover-report: ## an html report of the coverage statistics
	go test -v ./... -covermode=count -coverpkg=./... -coverprofile coverage.out
	go tool cover -html coverage.out -o coverage.html
	open coverage.html

lint: ## runs the linter
	golangci-lint -v run ./...

gen: ## runs go generate & sqlc commands
	go generate ./...
	docker run --rm -u $$(id -u ${USER}):$$(id -g ${USER}) -v $$(pwd):/src -w /src sqlc/sqlc generate -f postgres/sqlc.yaml

dc-up: ## starts the docker-compose environment
	docker compose -f docker-compose.yaml up -d

dc-down: ## stops the docker-compose environment
	docker compose -f docker-compose.yaml down

db-enter: ## enters the database container
	docker-compose -f docker-compose.yaml exec db psql -U postgres postgres

migrate-create: ## creates a migration file: Usage: make migrate-create name=<name>
	docker run --rm -u $$(id -u ${USER}):$$(id -g ${USER})  -v ${PWD}/scripts/migrations:/migrations migrate/migrate create -dir /migrations -ext sql -seq  ${name}

deploy: ## deploys the app
	fly deploy --ha=false
