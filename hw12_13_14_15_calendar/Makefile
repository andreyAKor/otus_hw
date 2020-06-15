PROJECTNAME=calendar

GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin

build:
	go build -ldflags="-s -w" -o '$(GOBIN)/$(PROJECTNAME)' ./cmd/$(PROJECTNAME)/main.go || exit

run:
	go build -o '$(GOBIN)/$(PROJECTNAME)' ./cmd/$(PROJECTNAME)/main.go
	'$(GOBIN)/$(PROJECTNAME)' --config='$(GOBASE)/configs/$(PROJECTNAME).yml'

test:
	go test -v -count=1 -race -timeout=60s ./...

install-deps:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint

lint: install-deps
	golangci-lint run ./...

install:
	go mod download

.PHONY: build
