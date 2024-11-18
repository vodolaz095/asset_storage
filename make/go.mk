deps:
	go mod download
	go mod verify
	go mod tidy

test:
	go test -v ./...

start:
	go run main.go

build:
	CGO_ENABLED=0 go build -o build/asset_storage main.go

.PHONY: build
