SHELL := /bin/bash
.DEFAULT_GOAL := help

# https://gist.github.com/tadashi-aikawa/da73d277a3c1ec6767ed48d1335900f3
.PHONY: $(shell grep -h -E '^[a-zA-Z_-]+:' $(MAKEFILE_LIST) | sed 's/://')

tools: ## install tools
	go install github.com/pilu/fresh github.com/gobuffalo/packr/v2/packr2

dev: ## start dev mode
	fresh

start: ## start server
	packr2
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -o app .
	packr2 clean
	docker-compose up -d

open: ## open browser
	open http://127.0.0.1:4455/

# https://postd.cc/auto-documented-makefile/
help: ## Show help message
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

