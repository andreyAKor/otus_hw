package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/repository/repository"

	"github.com/pkg/errors"
)

// Prepare event scructure from POST-data.
func (s Server) prepareEventRequest(r *http.Request) (*repository.Event, error) {
	date, err := time.Parse(time.RFC3339, r.PostFormValue("date"))
	if err != nil {
		return nil, errors.Wrap(err, `"date" parsing fail`)
	}

	duration, err := time.ParseDuration(r.PostFormValue("duration"))
	if err != nil {
		return nil, errors.Wrap(err, `"duration" parsing fail`)
	}

	userID, err := strconv.Atoi(r.PostFormValue("user_id"))
	if err != nil {
		return nil, errors.Wrap(err, `"user_id" parsing fail`)
	}

	ev := &repository.Event{
		Title:    r.PostFormValue("title"),
		Date:     date,
		Duration: duration,
		UserID:   int64(userID),
	}

	if descr := r.PostFormValue("descr"); len(descr) > 0 {
		ev.Descr = &descr
	}

	if durationStart := r.PostFormValue("duration_start"); len(durationStart) > 0 {
		durationStart, err := time.ParseDuration(durationStart)
		if err != nil {
			return nil, errors.Wrap(err, `"duration_start" parsing fail`)
		}

		ev.DurationStart = &durationStart
	}

	return ev, nil
}

// Preparing http-result events list from reository events list.
func (s Server) prepareEventsResponse(events []repository.Event) []Event {
	resEvents := make([]Event, len(events))
	for idx, ev := range events {
		resEv := Event{
			Title:    ev.Title,
			Date:     ev.Date,
			Duration: ev.Duration.String(),
			Descr:    ev.Descr,
			UserID:   ev.UserID,
		}

		if ev.DurationStart != nil {
			durationStart := ev.DurationStart.String()
			resEv.DurationStart = &durationStart
		}

		resEvents[idx] = resEv
	}

	return resEvents
}
