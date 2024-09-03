test-:
	go test ./...

.PHONY: restart
restart:
	docker compose restart

.PHONY: log
log:
	docker logs -f app

.PHONY: down

down:
	docker compose down

.PHONY: up
up:
	if [ "$(MODE)" = "l" ]; then \
		docker compose up; \
	else \
		docker compose up -d; \
	fi

.PHONY: psql
psql:
	docker exec -it db psql -U postgres teste
