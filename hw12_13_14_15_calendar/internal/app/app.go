package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/grpc"
	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/http"

	"github.com/pkg/errors"
)

type App struct {
	httpSrv *http.Server
	grpcSrv *grpc.Server
}

func New(httpSrv *http.Server, grpcSrv *grpc.Server) (*App, error) {
	return &App{httpSrv, grpcSrv}, nil
}

// Run application.
func (a *App) Run(ctx context.Context) error {
	errCh := make(chan error)

	// Init http-server
	go func() {
		if err := a.httpSrv.Run(ctx); err != nil {
			errCh <- errors.Wrap(err, "http-server listen fail")
		}
	}()

	// Init grpc-server
	go func() {
		if err := a.grpcSrv.Run(ctx); err != nil {
			errCh <- errors.Wrap(err, "grpc-server listen fail")
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
