default: clean test build run

test:
	go test ./...

build:
	go build -o ./build/app ./cmd/main.go

run:
	./build/app ${ARGS}

clean:
	rm -rf build
