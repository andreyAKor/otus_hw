package grpc

import (
	"time"

	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/repository/repository"
	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/schema"

	"github.com/golang/protobuf/ptypes/duration"
	"github.com/golang/protobuf/ptypes/timestamp"
)

// Preparing grpc request event to reository event representation.
func (s Server) prepareEventRequest(event *schema.Event) repository.Event {
	ev := repository.Event{
		Title:  event.Title,
		UserID: event.UserID,
	}

	if v := event.Date; v != nil {
		ev.Date = time.Unix(v.Seconds, int64(v.Nanos))
	}
	if v := event.Duration; v != nil {
		ev.Duration = time.Duration(v.Seconds*1e9 + int64(v.Nanos))
	}
	if len(event.Descr) > 0 {
		ev.Descr = &event.Descr
	}
	if v := event.DurationStart; v != nil {
		durationStart := time.Duration(v.Seconds*1e9 + int64(v.Nanos))
		ev.DurationStart = &durationStart
	}

	return ev
}

// Preparing grpc result events list from reository events list.
func (s Server) prepareEventsResponse(events []repository.Event) []*schema.Event {
	resEvents := make([]*schema.Event, len(events))
	for idx, ev := range events {
		resEv := &schema.Event{
			Title: ev.Title,
			Date: &timestamp.Timestamp{
				Seconds: ev.Date.Unix(),
				Nanos:   int32(ev.Date.Nanosecond()),
			},
			Duration: &duration.Duration{
				Seconds: int64(ev.Duration.Seconds()),
				Nanos:   int32(ev.Duration % time.Second),
			},
			UserID: ev.UserID,
		}

		if ev.Descr != nil {
			resEv.Descr = *ev.Descr
		}
		if ev.DurationStart != nil {
			resEv.DurationStart = &duration.Duration{
				Seconds: int64(ev.DurationStart.Seconds()),
				Nanos:   int32(*ev.DurationStart % time.Second),
			}
		}

		resEvents[idx] = resEv
	}

	return resEvents
}
