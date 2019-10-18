SERVER_OUT := "app"
SERVER_PKG_BUILD := "github.com/gidyon/app-beston"


build_server_prod: ## Build a production binary for server
	CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -a -installsuffix cgo -ldflags '-s' -v -o $(SERVER_OUT) $(SERVER_PKG_BUILD)
	
build_server: ## Build the binary file for server
	@go build -i -v -o $(SERVER_OUT) $(SERVER_PKG_BUILD)

clean_server: ## Remove server binary
	@rm -f $(SERVER_OUT)

docker_build: ## Create a docker image for the service
ifdef tag
	@docker build --no-cache -t gidyon/app-beston:$(tag) .
else
	@docker build --no-cache -t gidyon/app-beston:latest .
endif

docker_build_prod: build_server_prod docker_build

docker_tag:
	@docker tag gidyon/app-beston:$(tag) gidyon/gidyon/app-beston:$(tag)

docker_push:
	@docker push gidyon/gidyon/app-beston:$(tag)

docker_build_and_push: docker_build_prod docker_tag docker_push 
	
help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
