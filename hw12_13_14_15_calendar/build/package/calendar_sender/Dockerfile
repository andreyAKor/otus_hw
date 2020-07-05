# Environment
FROM golang:1.14 as build-env

RUN mkdir -p /opt/calendar
WORKDIR /opt/calendar
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /opt/bin/calendar_sender ./cmd/calendar_sender/main.go

# Release
FROM alpine:latest

COPY --from=build-env /opt/calendar/configs/calendar_sender.yml /etc/calendar_sender.yml
COPY --from=build-env /opt/bin/calendar_sender /bin/calendar_sender

RUN mkdir -p /var/log/calendar
ENV LOGGING_FILE=/var/log/calendar_sender/calendar_sender.log

ENTRYPOINT ["/bin/calendar_sender", "--config=/etc/calendar_sender.yml"]
