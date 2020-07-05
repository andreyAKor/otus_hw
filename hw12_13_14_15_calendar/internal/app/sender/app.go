package sender

import (
	"context"
	"io"

	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/rmq/consumer"

	"github.com/pkg/errors"
)

var _ io.Closer = (*App)(nil)

type App struct {
	cons *consumer.Consumer
}

func New(cons *consumer.Consumer) (*App, error) {
	return &App{cons}, nil
}

// Run application.
func (a *App) Run(ctx context.Context) error {
	// Init rmq-consumer
	if err := a.cons.Run(ctx); err != nil {
		return errors.Wrap(err, "rmq-consumer fail")
	}

	return nil
}

func (a *App) Close() error {
	return a.cons.Close()
}
