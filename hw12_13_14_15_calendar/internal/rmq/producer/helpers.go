package producer

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
)

func Worker(ctx context.Context, d time.Duration, fn func() error) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(d):
				if err := fn(); err != nil {
					log.Fatal().Err(err).Msg("worker error")
					return
				}
			}
		}
	}()
}
