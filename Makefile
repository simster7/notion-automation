
.PHONY: lint
lint:
	gofmt -w .

.PHONY: clean
clean:
	rm -rf dist

.PHONY: build
build:
	go build -o dist/main .
