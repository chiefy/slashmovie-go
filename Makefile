
include .env

vendor:
	@dep ensure

.PHONY: local-proxy
local-proxy:
	@ngrok http $(PORT)

.PHONY: run
run: vendor
	@PORT=$(PORT) TMDB_API_KEY=$(TMDB_API_KEY) SLACK_SIGNING_SECRET=$(SLACK_SIGNING_SECRET) go run .