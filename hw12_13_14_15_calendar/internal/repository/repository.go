package repository

import (
	"context"
	"strings"

	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/repository/memory"
	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/repository/psql"
	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/repository/repository"

	"github.com/pkg/errors"
)

var (
	ErrUnknowDatabaseType = errors.New("unknow database type")
)

// Init database.
func New(ctx context.Context, t, dsn string) (repository.EventsRepo, error) {
	// Init database type
	var r repository.EventsRepo

	switch strings.ToLower(t) {
	case "memory":
		r = new(memory.Repo)
	case "db":
		rsql := new(psql.Repo)
		if err := rsql.Connect(ctx, dsn); err != nil {
			return nil, errors.Wrap(err, "connection error")
		}
		r = rsql
	default:
		return nil, ErrUnknowDatabaseType
	}

	return r, nil
}
