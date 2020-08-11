package sender

import (
	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/configs"

	"github.com/pkg/errors"
)

var _ configs.Configer = (*Config)(nil)

type Config struct {
	Logging  configs.Logging
	Database configs.Database
	RMQ      configs.RMQ

	Queues struct {
		Events  configs.Queue
		Senders configs.Queue
	}

	// Consumer settings
	Consumer struct {
		// Consumer tag name
		ConsumerTag string

		QosPrefetchCount int

		// Count threads for reading queue messages
		Threads int
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
