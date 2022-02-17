test:
	go test ./...

build:
	go build -o cmd\main .\cmd\main.go

run:
	.\cmd\main