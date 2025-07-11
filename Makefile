.PHONY: up down reset e2e-ci-up e2e-ci-down e2e-ci-ps

up:
	docker-compose -f deployment/local/docker-compose.yaml up -d

down:
	docker-compose -f deployment/local/docker-compose.yaml down

reset:
	docker-compose -f deployment/local/docker-compose.yaml down -v

e2e-ci-up: ssl/nginx.key
	docker compose -f docker-compose.ci.yaml up -d

ssl/nginx.key:
	openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout ssl/nginx.key -out ssl/nginx.crt -config ssl/openssl.conf -extensions v3_req

e2e-ci-ps:
	docker compose -f docker-compose.ci.yaml ps

e2e-ci-down:
	docker compose -f docker-compose.ci.yaml down
