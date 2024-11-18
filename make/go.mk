deps:
	go mod download
	go mod verify
	go mod tidy

test:
	go test -v ./...

start:
	go run main.go
