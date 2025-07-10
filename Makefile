.PHONY: up down reset e2e-up e2e-down e2e-ps

up:
	docker-compose -f docker-compose.local.yaml up -d

down:
	docker-compose -f docker-compose.local.yaml down

reset:
	docker-compose -f docker-compose.local.yaml down -v

e2e-up: ssl/nginx.key
	docker compose -f docker-compose.ci.yaml up -d

ssl/nginx.key:
	openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout ssl/nginx.key -out ssl/nginx.crt -config ssl/openssl.conf -extensions v3_req

e2e-ps:
	docker compose -f docker-compose.ci.yaml ps

e2e-down:
	docker compose -f docker-compose.ci.yaml down
