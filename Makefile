docs:
	cd src && \
	go install github.com/swaggo/swag/cmd/swag@latest && \
	swag init --parseDependency --parseInternal 

run:
	cd src && go run .
	
build:
	cd src && docker build -t user-service .

start-all:
	make build && docker-compose up -d

stop-all:
	docker-compose down

start-postgres:
	docker-compose up -d postgres

start-user-service:
	docker-compose up -d user-service
