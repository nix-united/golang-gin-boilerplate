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

coverage_dir = coverage/integration

test-integration-coverage: # Runs integration tests against a -cover build of the deploy image and prints coverage
	@rm -rf $(coverage_dir) && mkdir -p $(coverage_dir)
	INTEGRATION_COVERDIR=$(CURDIR)/$(coverage_dir) go test -count=1 ./test/integration/
	@go tool covdata textfmt -i=$(coverage_dir) -o coverage/integration.out
	@go tool cover -html=coverage/integration.out -o coverage/integration.html
	@go tool covdata percent -i=$(coverage_dir)
	@go tool cover -func=coverage/integration.out | tail -n 1
	@echo "HTML report: coverage/integration.html"
