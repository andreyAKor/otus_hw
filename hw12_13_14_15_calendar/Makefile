GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin

build:
	go build -ldflags="-s -w" -o '$(GOBIN)/calendar' ./cmd/calendar/main.go || exit
	go build -ldflags="-s -w" -o '$(GOBIN)/calendar_scheduler' ./cmd/calendar_scheduler/main.go || exit
	#go build -ldflags="-s -w" -o '$(GOBIN)/calendar_sender' ./cmd/calendar_sender/main.go || exit

run:
	go build -o '$(GOBIN)/calendar' ./cmd/calendar/main.go
	'$(GOBIN)/calendar' --config='$(GOBASE)/configs/calendar.yml'

test:
	go test -v -count=1 -race -timeout=60s ./...

install-deps:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint

lint: install-deps
	golangci-lint run ./...

install:
	go mod download

generate:
	go generate ./...

.PHONY: build
