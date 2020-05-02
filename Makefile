lint_docker_compose_file = "./development/golangci_lint/docker-compose.yml"

lint-build:
	docker-compose --file=$(lint_docker_compose_file) build

lint-check: lint-build
	docker-compose --file=$(lint_docker_compose_file) run gin-golinter golangci-lint run -v

lint-fix: lint-build
	docker-compose --file=$(lint_docker_compose_file) run gin-golinter golangci-lint run --fix