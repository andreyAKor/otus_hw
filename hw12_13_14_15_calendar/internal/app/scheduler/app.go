package scheduler

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/rmq"
	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/rmq/producer"

	"github.com/pkg/errors"
)

type App struct {
	prod *producer.Producer
}

func New(r *rmq.Rmq) (*App, error) {
	prod, err := producer.New(r)
	if err != nil {
		return nil, errors.Wrap(err, "producer init fail")
	}

	return &App{prod}, nil
}

// Run application.
func (a *App) Run(ctx context.Context) error {
	errCh := make(chan error)

	// Init rmq-producer
	go func() {
		if err := a.prod.Run(ctx); err != nil {
			errCh <- errors.Wrap(err, "rmq-producer fail")
		}
	}()

	// Graceful shutdown
	interruptCh := make(chan os.Signal, 1)
	signal.Notify(interruptCh, os.Interrupt, syscall.SIGTERM)

END:
	for {
		select {
		case <-ctx.Done():
			break END
		case <-interruptCh:
			break END
		case err := <-errCh:
			return err
		}
	}

	return nil
}
