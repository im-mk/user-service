run:
	cd src && \
	go install github.com/swaggo/swag/cmd/swag@latest && \
	swag init --parseDependency --parseInternal && \
	go run .
	
build:
	cd src && docker build -t user-service .

start-all:
	make build && \
	cd infra/local && make start

stop-all:
	cd infra/local && make stop
