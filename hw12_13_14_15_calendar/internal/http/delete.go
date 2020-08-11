package http

import (
	"net/http"
	"strconv"

	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/repository"

	"github.com/pkg/errors"
)

// Delete event.
func (s *Server) delete(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return nil, errors.Wrap(err, `"id" parsing fail`)
	}

	if err := s.calendar.Delete(r.Context(), int64(id)); err != nil {
		switch err {
		case repository.ErrNotFound:
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}

		return nil, errors.Wrap(err, "can't delete event")
	}

	return nil, nil
}
