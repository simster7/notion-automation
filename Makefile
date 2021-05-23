.DEFAULT_GOAL := build

.PHONY: build
build:
	go build -o dist/main .

.PHONY: lint
lint:
	gofmt -w .

.PHONY: clean
clean:
	rm -rf dist

