.PHONY:	help
help: ## show help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY:	run-server
run-server: ## go run main.go --server
	@go run ./main.go --server

.PHONY:	run
run: ## go run main.go
	@go run ./main.go
