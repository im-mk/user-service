build-user-service:
	cd src && docker build -t user-service .

start:
	make build-user-service && \
	cd infra/local && make start

stop:
	cd infra/local && make stop

run:
	cd src && \
	go install github.com/swaggo/swag/cmd/swag@latest && \
	swag init --parseDependency --parseInternal && \
	go run .