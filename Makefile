.PHONY: create-keys start-pgadmin start-user-db test-user-api build-user-api start-user-api start stop down

create-keys:
	mkdir -p src/.keys	
	openssl genrsa -out src/.keys/private.pem 2048
	openssl rsa -in src/.keys/private.pem -pubout -out src/.keys/public.pem

start-pgadmin:
	docker compose up -d user-service-pgadmin

start-user-db:
	docker compose up -d user-service-db user-service-liquibase

test-user-api:
	docker build -f src/Dockerfile --target test -t user-service-test ./src
	docker run --rm user-service-test

build-user-api:
	docker compose build user-service

start-user-api: create-keys test-user-api build-user-api start-user-db
	docker compose up -d user-service

start:
	docker compose up -d --build

stop:
	docker compose stop

down: 
	docker compose down -v

