package calendar

import (
	"context"
	"io"

	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/grpc"
	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/http"

	"github.com/rs/zerolog/log"
)

var _ io.Closer = (*App)(nil)

type App struct {
	httpSrv *http.Server
	grpcSrv *grpc.Server
}

func New(httpSrv *http.Server, grpcSrv *grpc.Server) (*App, error) {
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

func (a *App) Close() error {
	if err := a.httpSrv.Close(); err != nil {
		return err
	}
	if err := a.grpcSrv.Close(); err != nil {
		return err
	}

	return nil
}
