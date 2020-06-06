PROJECTNAME=calendar

GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin

build:
	go build -ldflags="-s -w" -o '$(GOBIN)/$(PROJECTNAME)' ./cmd/$(PROJECTNAME)/main.go || exit

run:
	go build -o '$(GOBIN)/$(PROJECTNAME)' ./cmd/$(PROJECTNAME)/main.go || exit
	'$(GOBIN)/$(PROJECTNAME)' --config='$(GOBASE)/configs/$(PROJECTNAME).yml'

test:
	go test -v -count=1 -race -gcflags=-l -timeout=30s ./...

lint:
	golangci-lint run ./...

install:
	go mod download

.PHONY: build
