GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin

build:
	go build -ldflags="-s -w" -o '$(GOBIN)/calendar' ./cmd/calendar/main.go || exit
	go build -ldflags="-s -w" -o '$(GOBIN)/calendar_scheduler' ./cmd/calendar_scheduler/main.go || exit
	go build -ldflags="-s -w" -o '$(GOBIN)/calendar_sender' ./cmd/calendar_sender/main.go || exit

run:
	go build -o '$(GOBIN)/calendar' ./cmd/calendar/main.go
	'$(GOBIN)/calendar' --config='$(GOBASE)/configs/calendar.yml'

up:
	docker-compose up -d --build

down:
	docker-compose down

test:
	go test -v -count=1 -race -timeout=60s ./...

integration-tests:
	set -e ;\
	docker-compose -f docker-compose.yml -f docker-compose.test.yml up --build -d ;\
	test_status_code=0 ;\
	docker-compose -f docker-compose.yml -f docker-compose.test.yml run integration_tests go test || test_status_code=$$? ;\
	docker-compose -f docker-compose.yml -f docker-compose.test.yml down ;\
	exit $$test_status_code ;

install-deps:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint

lint: install-deps
	golangci-lint run ./...

install:
	go mod download

generate:
	go generate ./...

.PHONY: build
