# Load environment variables from .env file
include .env

build: templ
	docker-compose build

templ:
	templ generate

# Define the start task
start: build
	docker-compose up -d

# Define the test task, should call the container test from the docker compose
test:
	@if command -v jq >/dev/null; then \
		curl $(GRAFANA_USER_URL) | jq . ; \
	else \
		curl $(GRAFANA_USER_URL); \
	fi

health:
	curl $(HEALTH_URL);

# Define the stop task
stop:
	docker-compose down

# Define the clean task
clean: stop
	docker-compose rm -f
	rm -f $(LOCAL_ASSUME_FILE_LOCATION)

# Define the logs task
logs:
	docker-compose logs -f

# Define the shell task
shell:
	docker-compose exec app /bin/bash

# Define the run task that calls clean, start, waits for 30 seconds and then calls test
run: clean start
	sleep 5
	make test
