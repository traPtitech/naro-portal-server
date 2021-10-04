.PHONY: help
help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

.PHONY: init
init: ## Download and install go mod dependencies
	@go mod download

.PHONY: lint
lint: ## Lint go files
	@golangci-lint run

.PHONY: build
build: ## Build binary
	@go build -o Q-n-A ./*.go

.PHONY: run
run: ## Run Q'n'A directly
	@go run ./*.go

.PHONY: up
up: ## Build and start Q'n'A debug containers (Not recommended)
	@docker-compose up -d --force-recreate

.PHONY: down
down: ## Stop and remove app containers (Not recommended as `up`)
	@docker-compose down

.PHONY: reset-frontend
reset-frontend: stop-front rm-front delete-front-image ## Delete frontend container and image completely

.PHONY: stop-front
stop-front:
	@docker ps -a | grep Q-n-A_frontend | awk '{print $$1}' | xargs docker stop

.PHONY: rm-front
rm-front:
	@docker ps -a | grep Q-n-A_frontend | awk '{print $$1}' | xargs docker rm

.PHONY: delete-front-image
delete-front-image:
	@docker images -a | grep q-n-a | grep frontend | awk '{print $$3}' | xargs docker rmi

.PHONY: chown
chown:
	$(eval name := $(shell whoami))
	@sudo chown -R $(name):$(name) .