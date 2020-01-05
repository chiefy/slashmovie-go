include .env

VERSION:=1.0.0

project:=$(shell basename $(shell pwd))
commit:=$(shell git rev-parse --short HEAD)
importpath:=github.com/chiefy/$(project)
ts:=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
binary:=slashmovie 

$(GOPATH)/bin/dep:
	@curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

vendor: $(GOPATH)/bin/dep
	@dep ensure

.PHONY: build
build: $(binary)

$(binary): vendor
	go build -ldflags \
	"-X main.Version=$(VERSION) \
	-X main.Commit=$(commit) \
	-X main.Date=$(ts)" \
	-o $(binary) .

.PHONY: local-proxy
local-proxy:
	@ngrok http $(PORT)

.PHONY: run
run: vendor
	@PORT=$(PORT) TMDB_API_KEY=$(TMDB_API_KEY) SLACK_SIGNING_SECRET=$(SLACK_SIGNING_SECRET) ./slashmovie