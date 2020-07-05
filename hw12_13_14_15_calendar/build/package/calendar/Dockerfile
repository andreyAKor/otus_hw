# Environment
FROM golang:1.14 as build-env

RUN mkdir -p /opt/calendar
WORKDIR /opt/calendar
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /opt/bin/calendar ./cmd/calendar/main.go

# Release
FROM alpine:latest

COPY --from=build-env /opt/calendar/configs/calendar.yml /etc/calendar.yml
COPY --from=build-env /opt/bin/calendar /bin/calendar

RUN mkdir -p /var/log/calendar
ENV LOGGING_FILE=/var/log/calendar/calendar.log

ENTRYPOINT ["/bin/calendar", "--config=/etc/calendar.yml"]
