.PHONY: build clean test

hash = $(shell git log --pretty=format:'%h' -n 1)

build: clean
	mkdir -p build
	go build -o build/narrative -ldflags "-X main.build=${hash}" cmd/narrative/main.go

test:
	go test ./...

run:
	go run cmd/narrative/main.go

clean:
	rm -rf build
