build:
	golangci-lint run
	go test ./...

unit-tests:
	go test -short ./...

lint:
	golangci-lint run

run:
	go run cmd/httpserver/main.go -endpoints=specifications/examples/ -dev-mode=true