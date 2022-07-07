all: swagger_gen test run
swagger_gen:
#Run swag init to create documents
	go install github.com/swaggo/swag/cmd/swag@v1.8.3
	swag init

test:
	go test ./...

run:
	docker compose up -d