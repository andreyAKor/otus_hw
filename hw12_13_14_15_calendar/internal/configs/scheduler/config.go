package scheduler

import (
	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/configs"

	"github.com/pkg/errors"
)

var _ configs.Configer = (*Config)(nil)

type Config struct {
	Logging  configs.Logging
	Database configs.Database

	// RabbitMQ settings
	RMQ struct {
		// URI connection string to RabbitMQ
		URI string

		// Exchange name
		ExchangeName string

		// Exchange type
		ExchangeType string

		// Queue name
		QueueName string

		// Message routing rules
		BindingKey string
	}
}

func (c *Config) Init(file string) error {
	cfg, err := configs.Init(file, c)

	c, ok := cfg.(*Config)
	if !ok {
		return errors.Wrap(err, "init config failed")
	}

	return nil
}
