# Environment
FROM golang:1.14

RUN mkdir -p /opt/goose
WORKDIR /opt/goose

RUN go get -u github.com/pressly/goose/cmd/goose

COPY . .

ENTRYPOINT ["/go/bin/goose", "-dir", "migrations", "up"]
