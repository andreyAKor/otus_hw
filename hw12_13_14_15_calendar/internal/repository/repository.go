package repository

import (
	"context"
	"time"

	"github.com/pkg/errors"
)

var (
	ErrTimeBusy = errors.New("event at this time is busy")
	ErrNotFound = errors.New("event not found")
)

//go:generate mockgen -source=$GOFILE -destination ./mocks/mock_events_repo.go -package mocks EventsRepo
type EventsRepo interface {
	Create(ctx context.Context, ev Event) (int64, error)
	Update(ctx context.Context, id int64, ev Event) error

	Delete(ctx context.Context, id int64) error
	DeleteOld(ctx context.Context) error

	GetListByDate(ctx context.Context, date time.Time) ([]Event, error)
	GetListByWeek(ctx context.Context, start time.Time) ([]Event, error)
	GetListByMonth(ctx context.Context, start time.Time) ([]Event, error)
}

type DBEventsRepo interface {
	Connect(ctx context.Context, dsn string) error
	Close() error

	EventsRepo
}

type Event struct {
	ID            int64
	Title         string
	Date          time.Time
	Duration      time.Duration
	Descr         *string
	UserID        int64
	DurationStart *time.Duration
	CreatedAt     time.Time
	UpdatedAt     *time.Time
}
