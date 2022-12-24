-include .env

VERSION:=2.1.0
DOCKER_REPO_USER=chiefy
DOCKER_IMAGE_NAME=$(DOCKER_REPO_USER)/$(binary):$(VERSION)

project:=$(shell basename $(shell pwd))
commit:=$(shell git rev-parse --short HEAD)
importpath:=github.com/chiefy/$(project)
ts:=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
binary:=slashmovie

.PHONY: build
build: $(binary)

$(binary):
	go build \
	-tags netgo \
	-ldflags \
	"-s \
	-w \
	-X main.Version=$(VERSION) \
	-X main.Commit=$(commit) \
	-X main.Date=$(ts)" \
	-o $(binary) .

.PHONY: local-proxy
local-proxy:
	@ngrok http $(PORT)

.PHONY: clean 
clean:
	@-rm -f $(binary)
	
.PHONY: run
run: $(binary)
	@./$(binary)

.PHONY: build-docker
build-docker:
	docker build -t $(DOCKER_IMAGE_NAME) .

.PHONY: push-docker
push-docker:
	@docker push $(DOCKER_IMAGE_NAME)
