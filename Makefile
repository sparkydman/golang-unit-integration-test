build: build
	@go build -o bin/efficient_api

run: build
	@./bin/efficient_api

test:
	@go test ./... -v

testrace:
	@go test ./... -v --race

testcover:
	@go test ./... -v -cover