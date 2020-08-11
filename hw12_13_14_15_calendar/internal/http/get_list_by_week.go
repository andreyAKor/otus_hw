package http

import (
	"net/http"
	"time"

	"github.com/pkg/errors"
)

type getListByWeek struct {
	EventList
}

// Get events list by week.
func (s *Server) getListByWeek(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	start, err := time.Parse("2006-01-02", r.URL.Query().Get("start"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return nil, errors.Wrap(err, `"start" parsing fail`)
	}

	events, err := s.calendar.GetListByWeek(r.Context(), start)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return nil, errors.Wrap(err, "can't get events list by week")
	}

	resEvents := s.prepareEventsResponse(events)
	return getListByWeek{
		EventList{
			Events: resEvents,
		},
	}, nil
}
