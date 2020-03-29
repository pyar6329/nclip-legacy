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

.PHONY:	build-darwin
build-darwin: ## go build to darwin amd64
	@CGO_ENABLED=1 GOARCH="amd64" GOOS="darwin" go build -a -tags netgo -ldflags '-w -s -extldflags "-static"' -o build/darwin/nclip *.go

.PHONY:	build-linux
build-linux: ## go build to linux amd64
	@CGO_ENABLED=1 GOARCH="amd64" GOOS="linux" go build -a -tags netgo -ldflags '-w -s -extldflags "-static"' -o build/linux/nclip *.go

.PHONY:	clean
clean: ## clean build binary
	@rm -rf ./build

.PHONY:	build
build: ## go build all binary
	@make build-darwin
	@make build-linux

.PHONY:	build-docker
build-docker: ## go build on docker
	@docker build -t foo:1 .
	@docker run --rm -v $(PWD)/build:/tmp/build foo:1 bash -c "cp -rf /build/* /tmp/build/"
