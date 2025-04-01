.PHONY: build
build:
	go build -o bin/instrument ./cmd/instrument

.PHONY: test
test:
	go test -v ./...

.PHONY: clean
clean:
	rm -rf bin/

.DEFAULT_GOAL := build 