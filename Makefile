.PHONY: docs test run build start stop start-postgres start-user-service

docs:
	cd src && \
	go run github.com/swaggo/swag/cmd/swag@latest init --parseDependency --parseInternal

test:
	cd src && go test ./... 

run:
	cd src && go run .
	
build:
	cd src && go test ./... -v && docker build -t user-service .

start:
	make build && docker compose up -d

stop:
	docker compose down

start-postgres:
	docker compose up -d postgres

start-user-service:
	docker compose up -d user-service
