package http

import (
	"net/http"
	"time"

	"github.com/pkg/errors"
)

type getListByDate struct {
	EventList
}

// Get events list by date.
func (s *Server) getListByDate(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	date, err := time.Parse("2006-01-02", r.URL.Query().Get("date"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return nil, errors.Wrap(err, `"date" parsing fail`)
	}

	events, err := s.calendar.GetListByDate(r.Context(), date)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return nil, errors.Wrap(err, "can't get events list by date")
	}

	resEvents := s.prepareEventsResponse(events)
	return getListByDate{
		EventList{
			Events: resEvents,
		},
	}, nil
}
