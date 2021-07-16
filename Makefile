.DEFAULT_GOAL := build
TOKEN := $(shell cat secrets/notion_secret.txt)
CALENDAR_ID := $(shell cat secrets/calendar_id.txt)

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
	gcloud functions deploy Nightly --runtime go113 --trigger-http --allow-unauthenticated --set-env-vars "NOTION_TOKEN=$(TOKEN)" --set-env-vars "CALENDAR_ID=$(CALENDAR_ID)"
