package http

import (
	"net/http"
	"strconv"

	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/repository"

	"github.com/pkg/errors"
)

// Updating event.
func (s *Server) update(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	ev, err := s.prepareEventRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return nil, errors.Wrap(err, "can't preparing event")
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return nil, errors.Wrap(err, `"id" parsing fail: `+r.URL.Query().Get("id"))
	}

	if err := s.r.Update(r.Context(), int64(id), *ev); err != nil {
		switch err {
		case repository.ErrTimeBusy:
			w.WriteHeader(http.StatusConflict)
		case repository.ErrNotFound:
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}

		return nil, errors.Wrap(err, "can't update event")
	}

	return nil, nil
}
