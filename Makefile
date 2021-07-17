.DEFAULT_GOAL := build
TOKEN := $(shell cat secrets/notion_secret.txt)
CALENDAR_ID := $(shell cat secrets/calendar_id.txt)
SERVICE_ACCOUNT := $(shell cat secrets/service_account_name.txt)

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
	gcloud functions deploy Nightly --runtime go113 --trigger-http --allow-unauthenticated --service-account "$(SERVICE_ACCOUNT)" --set-env-vars "NOTION_TOKEN=$(TOKEN)" --set-env-vars "CALENDAR_ID=$(CALENDAR_ID)"
