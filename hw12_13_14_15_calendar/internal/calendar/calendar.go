package calendar

import (
	"context"
	"io"
	"strings"
	"time"

	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/repository"
	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/repository/memory"
	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/repository/psql"

	"github.com/pkg/errors"
)

//go:generate mockgen -source=$GOFILE -destination ./mocks/mock_calendarer.go -package mocks Calendarer
type Calendarer interface {
	repository.EventsRepo
}

var (
	ErrUnknowDatabaseType = errors.New("unknow database type")

	_ io.Closer  = (*Calendar)(nil)
	_ Calendarer = (*Calendar)(nil)
)

type Calendar struct {
	repo repository.EventsRepo
}

func New(ctx context.Context, t, dsn string) (*Calendar, error) {
	// Init database type
	var repo repository.EventsRepo

	switch strings.ToLower(t) {
	case "memory":
		repo = new(memory.Repo)
	case "db":
		rsql := new(psql.Repo)
		if err := rsql.Connect(ctx, dsn); err != nil {
			return nil, errors.Wrap(err, "connection error")
		}
		repo = rsql
	default:
		return nil, ErrUnknowDatabaseType
	}

	return &Calendar{repo}, nil
}

func (c *Calendar) Close() error {
	if v, ok := c.repo.(repository.DBEventsRepo); ok {
		return v.Close()
	}

	return nil
}

func (c *Calendar) Create(ctx context.Context, ev repository.Event) (int64, error) {
	return c.repo.Create(ctx, ev)
}

func (c *Calendar) Update(ctx context.Context, id int64, ev repository.Event) error {
	return c.repo.Update(ctx, id, ev)
}

func (c *Calendar) Delete(ctx context.Context, id int64) error {
	return c.repo.Delete(ctx, id)
}

func (c *Calendar) DeleteOld(ctx context.Context) error {
	return c.repo.DeleteOld(ctx)
}

func (c *Calendar) GetListByDate(ctx context.Context, date time.Time) ([]repository.Event, error) {
	return c.repo.GetListByDate(ctx, date)
}

func (c *Calendar) GetListByWeek(ctx context.Context, start time.Time) ([]repository.Event, error) {
	return c.repo.GetListByWeek(ctx, start)
}

func (c *Calendar) GetListByMonth(ctx context.Context, start time.Time) ([]repository.Event, error) {
	return c.repo.GetListByMonth(ctx, start)
}
