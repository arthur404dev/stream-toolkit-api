.PHONY= all deps lint fmt
.EXPORT_ALL_VARIABLES:

APP_NAME=api
APP_VERSION=1.0
ENTRY_FILES=main.go 

ifneq (,$(wildcard ./.env))
    include .env
    export
endif

all: 
	@echo "=>> This is the app: $(APP_NAME) using $(ENVIRONMENT)"
deps:
	@echo "=>> Updating Dependencies"
	@go get -u=patch
lint:
	@echo "=>> Linting code with vet"
	@go vet ./...
fmt:
	@echo "=>> Formatting code"
	@go fmt ./...
run-dev:
	@echo "=>> Running locally with air"
	@air run $(ENTRY_FILES)
docker-build:
	@echo "=>> Building Docker Image"
	@docker build --tag $(APP_NAME):$(APP_VERSION) .
run-prod:
	@echo "=>> Running Production release"
	docker run -d --rm -ti --name=$(APP_NAME) -p $(PORT):$(PORT) $(APP_NAME):$(APP_VERSION)