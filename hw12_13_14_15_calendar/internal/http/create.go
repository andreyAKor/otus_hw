package http

import (
	"net/http"

	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/repository"

	"github.com/pkg/errors"
)

type create struct {
	ID int64 `json:"id"`
}

// Add new event.
func (s *Server) create(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	ev, err := s.prepareEventRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return nil, errors.Wrap(err, "can't preparing event")
	}

	id, err := s.r.Create(r.Context(), *ev)
	if err != nil {
		switch err {
		case repository.ErrTimeBusy:
			w.WriteHeader(http.StatusConflict)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}

		return nil, errors.Wrap(err, "can't create new event")
	}

	w.WriteHeader(http.StatusCreated)
	return create{id}, nil
}
