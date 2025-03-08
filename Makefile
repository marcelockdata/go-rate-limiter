include .env
export



## ---------- UTILS
.PHONY: help
help: ## Show this menu
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: clean
clean: ## Clean all temp files
	@sudo rm -rf coverage*


## ---------- MAIN
.PHONY: test
test: ## Run unit-tests
	@go test -v ./... -coverprofile=coverage.out
	@go tool cover -html coverage.out -o coverage.html

.PHONY: build
build: ## Build the container image
	@docker build -t aleroxac/goexpert-ratelimiter:dev -f Dockerfile .

.PHONY: up
up: ## Put the compose containers up
	@docker-compose up -d

.PHONY: down
down: ## Put the compose containers down
	@docker-compose down



## ---------- SCENARIOS
define scenario_ip
	echo -n "Running IP rate limit scenario..."; \
	for i in {1..4}; do \
		echo -en "\nRequest $$i: "; \
		curl -is -w "%{http_code}" -o /dev/null http://localhost:8080/api/v1/zipcode/01001001; \
	done; \
	echo -e "\nWait for block duration: $(BLOCK_DURATION)s"; \
	sleep $(BLOCK_DURATION); \
	echo -en "Request after block: "; \
	curl -is -w "%{http_code}" -o /dev/null http://localhost:8080/api/v1/zipcode/01001001
endef

define scenario_token
	echo -n "Running token rate limit scenario..."; \
	for i in {1..6}; do \
		echo -en "\nRequest $$i: "; \
		curl -is -w "%{http_code}" -o /dev/null -H "API_KEY: my-token" http://localhost:8080/api/v1/zipcode/01001001; \
	done; \
	echo -e "\nWait for block duration: $(BLOCK_DURATION)s"; \
	sleep $(BLOCK_DURATION); \
	echo -en "Request after block: "; \
	curl -is -w "%{http_code}" -o /dev/null -H "API_KEY: my-token" http://localhost:8080/api/v1/zipcode/01001001
endef

.PHONY: run
run: ## Run test scenarios
	@if [ "$(SCENARIO)" = "ip" ]; then \
		$(call scenario_ip); \
	elif [ "$(SCENARIO)" = "token" ]; then \
		$(call scenario_token); \
	elif [ "$(SCENARIO)" = "all" ]; then \
		$(call scenario_ip); \
		echo -e "\n----------------------------------------"; \
		$(call scenario_token); \
	else \
		echo "Please specify a valid SCENARIO: (ip, token, all)"; \
	fi

