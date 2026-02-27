.PHONY: up down reset test run e2e-ci-up e2e-ci-down e2e-ci-ps

up: deployment/local/.env
	docker-compose -f deployment/local/docker-compose.yaml up -d

deployment/local/.env:
	cp .env deployment/local/.env

down:
	docker-compose -f deployment/local/docker-compose.yaml down

reset:
	docker-compose -f deployment/local/docker-compose.yaml down -v

test:
	./runtests.md	

run:
	go run main.go

e2e-ci-build:
	docker compose -f deployment/ci/docker-compose.yaml build

e2e-ci-up: deployment/ci/ssl/nginx.key
	docker compose -f deployment/ci/docker-compose.yaml up -d

deployment/ci/ssl/nginx.key:
	openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout deployment/ci/ssl/nginx.key -out deployment/ci/ssl/nginx.crt -config deployment/ci/ssl/openssl.conf -extensions v3_req

e2e-ci-ps:
	docker compose -f deployment/ci/docker-compose.yaml ps

e2e-ci-down:
	docker compose -f deployment/ci/docker-compose.yaml down

e2e-debug-ci-up: deployment/ci/ssl/nginx.key
	docker compose -f deployment/local/docker-compose.debug-ci.yaml up -d

e2e-debug-ci-ps:
	docker compose -f deployment/local/docker-compose.debug-ci.yaml ps

e2e-debug-ci-down:
	docker compose -f deployment/local/docker-compose.debug-ci.yaml down