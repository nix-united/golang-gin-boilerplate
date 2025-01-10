lint_docker_compose_file = "./development/golangci_lint/compose.yml"

lint-build:
	@echo "🌀 ️container are building..."
	@docker-compose --file=$(lint_docker_compose_file) build -q
	@echo "✔  ️container built"

lint-check:
	@echo "🌀️ code linting..."
	@docker-compose --file=$(lint_docker_compose_file) run --rm gin-golinter golangci-lint run \
		&& echo "✔️  checked without errors" \
		|| echo "☢️  code style issues found"

lint-fix:
	@echo "🌀 ️code fixing..."
	@docker-compose --file=$(lint_docker_compose_file) run --rm gin-golinter golangci-lint run --fix \
		&& echo "✔️  fixed without errors" \
		|| (echo "⚠️️  you need to fix above issues manually" && exit 1)
	@echo "⚠️️ run \"make lint-check\" again to check what did not fix yet"

organize-imports:
	@gci write --custom-order -s standard -s "prefix(github.com/nix-united/golang-gin-boilerplate)" -s default --skip-generated --skip-vendor .
