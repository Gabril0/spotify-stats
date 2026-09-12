run:
	go run cmd/server/main.go

swagger:
	swag init -g cmd/server/main.go --output docs
