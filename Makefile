build:
	@go build -o bin/iam-server cmd/server/main.go

run: build
	@./bin/iam-server

deps:
	@go mod download
	@go mod verify

tidy:
	@go mod tidy