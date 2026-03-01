.PHONY: docs test start

docs:
	go run github.com/swaggo/swag/cmd/swag@latest init --parseDependency --parseInternal

test:
	go test ./...

build:
	go build -o user-service .
	
start:
	go run .
