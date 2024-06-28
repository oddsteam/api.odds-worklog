up:
	docker-compose -f docker-compose.local.yaml up -d

down:
	docker-compose -f docker-compose.local.yaml down

reset:
	docker-compose -f docker-compose.local.yaml down -v
