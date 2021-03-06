package scheduler

import (
	"context"

	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/rmq/producer"

	"github.com/pkg/errors"
)

type App struct {
	prod producer.Producer
}

func New(prod producer.Producer) (*App, error) {
	return &App{prod}, nil
}

// Run application.
func (a *App) Run(ctx context.Context) error {
	// Init rmq-producer
	if err := a.prod.Run(ctx); err != nil {
		return errors.Wrap(err, "rmq-producer fail")
	}

	return nil
}

func (a *App) Close() error {
	return a.prod.Close()
}
