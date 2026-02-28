.PHONY: docs test start

docs:
	go run github.com/swaggo/swag/cmd/swag@latest init --parseDependency --parseInternal

test:
	go test ./... 

start:
	go run .
