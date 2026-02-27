.PHONY: up down reset test run

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
