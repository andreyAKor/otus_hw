package calendar

import (
	"context"

	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/calendar"
	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/grpc"
	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/http"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type App struct {
	httpSrv *http.Server
	grpcSrv *grpc.Server
}

func New(
	calendar calendar.Calendarer,
	httpHost string, httpPort int,
	grpcHost string, grpcPort int,
) (*App, error) {
	// Init http-server
	httpSrv, err := http.New(calendar, httpHost, httpPort)
	if err != nil {
		return nil, errors.Wrap(err, "can't initialize http-server")
	}

	// Init grpc-server
	grpcSrv, err := grpc.New(calendar, grpcHost, grpcPort)
	if err != nil {
		return nil, errors.Wrap(err, "can't initialize grpc-server")
	}

	return &App{httpSrv, grpcSrv}, nil
}

// Run application.
func (a *App) Run(ctx context.Context) error {
	go func() {
		if err := a.httpSrv.Run(ctx); err != nil {
			log.Fatal().Err(err).Msg("http-server listen fail")
		}
	}()
	go func() {
		if err := a.grpcSrv.Run(ctx); err != nil {
			log.Fatal().Err(err).Msg("grpc-server listen fail")
		}
	}()

	return nil
}

func (a *App) Close() (result error) {
	if err := a.httpSrv.Close(); err != nil {
		result = multierror.Append(result, err)
	}
	if err := a.grpcSrv.Close(); err != nil {
		result = multierror.Append(result, err)
	}

	return
}
