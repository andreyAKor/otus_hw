package http

import (
	"net/http"
	"time"

	"github.com/pkg/errors"
)

type getListByMonth struct {
	EventList
}

// Get events list by month.
func (s *Server) getListByMonth(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	start, err := time.Parse("2006-01-02", r.URL.Query().Get("start"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return nil, errors.Wrap(err, `"start" parsing fail`)
	}

	events, err := s.r.GetListByMonth(r.Context(), start)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return nil, errors.Wrap(err, "can't get events list by month")
	}

	resEvents := s.prepareEventsResponse(events)
	return getListByMonth{
		EventList{
			Events: resEvents,
		},
	}, nil
}
