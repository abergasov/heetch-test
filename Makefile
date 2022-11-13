
up:
	docker compose build
	docker compose up --build

lint: ## Check code style
	golangci-lint run --verbose

test: ## Run tests
	go test -race ./...