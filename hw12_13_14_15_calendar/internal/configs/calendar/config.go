package calendar

import (
	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/configs"

	"github.com/pkg/errors"
)

var _ configs.Configer = (*Config)(nil)

type Config struct {
	Logging  configs.Logging
	Database configs.Database

	// HTTP-server settings
	HTTP struct {
		// Host
		Host string

		// Port
		Port int
	}

	// GRPC-server settings
	GRPC struct {
		// Host
		Host string

		// Port
		Port int
	}
}

func (c *Config) Init(file string) error {
	cfg, err := configs.Init(file, c)

	_, ok := cfg.(*Config)
	if !ok {
		return errors.Wrap(err, "init config failed")
	}

	return nil
}
