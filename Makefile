DOCKER_BUILDKIT := 1

.PHONY:	help
help: ## show help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY:	run-server
run-server: ## go run main.go --server
	@go run ./main.go --server

.PHONY:	run
run: ## go run main.go
	@go run ./main.go

.PHONY:	build-darwin-intel
build-darwin-intel: ## go build to darwin intel CPU
	@CGO_ENABLED=1 GOARCH="amd64" GOOS="darwin" go build -a -tags netgo -ldflags '-w -s -extldflags "-static"' -o build/darwin/amd64/nclip *.go

.PHONY:	build-darwin-arm
build-darwin-arm: ## go build to darwin ARM CPU
	@CGO_ENABLED=1 GOARCH="arm64" GOOS="darwin" go build -a -tags netgo -ldflags '-w -s -extldflags "-static"' -o build/darwin/arm64/nclip *.go

.PHONY:	build-linux
build-linux: ## go build to linux amd64
	@CGO_ENABLED=1 GOARCH="amd64" GOOS="linux" go build -a -tags netgo -ldflags '-w -s -extldflags "-static"' -o build/linux/nclip *.go

.PHONY:	clean
clean: ## clean build binary
	@rm -rf ./build

.PHONY:	build
build: ## go build all binary
	@make build-darwin-intel
	@make build-darwin-arm
	@make build-linux
