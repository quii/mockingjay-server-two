build:
	golangci-lint run
	go test ./...

unit-tests:
	go test -short ./...

lint:
	golangci-lint run

run:
	go run cmd/mockingjay/main.go -endpoints=specifications/examples/ -dev-mode=true

run-cdcs:
	go run cmd/mockingjay/main.go -endpoints=specifications/examples/ -dev-mode=true -cdc=true