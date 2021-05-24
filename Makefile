.DEFAULT_GOAL := build
TOKEN := $(shell cat notion_secret.txt)

.PHONY: build
build:
	go build -o dist/main .

.PHONY: lint
lint:
	gofmt -w .

.PHONY: clean
clean:
	rm -rf dist

.PHONY: deploy
deploy:
	gcloud functions deploy Nightly --runtime go113 --trigger-http --allow-unauthenticated --set-env-vars "NOTION_TOKEN=$(TOKEN)"
