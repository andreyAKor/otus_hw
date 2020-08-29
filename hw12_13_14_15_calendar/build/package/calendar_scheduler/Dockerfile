# Environment
FROM golang:1.14 as build-env

RUN mkdir -p /opt/calendar
WORKDIR /opt/calendar
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /opt/bin/calendar_scheduler ./cmd/calendar_scheduler/main.go

# Release
FROM alpine:latest

COPY --from=build-env /opt/calendar/configs/calendar_scheduler.yml /etc/calendar_scheduler.yml
COPY --from=build-env /opt/bin/calendar_scheduler /bin/calendar_scheduler

RUN mkdir -p /var/log/calendar
ENV LOGGING_FILE=/var/log/calendar_scheduler/calendar_scheduler.log

ENTRYPOINT ["/bin/calendar_scheduler", "--config=/etc/calendar_scheduler.yml"]
