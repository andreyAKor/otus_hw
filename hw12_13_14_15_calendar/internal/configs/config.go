package configs

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Configer interface {
	Init(file string) error
}

// Logging settings.
type Logging struct {
	// Path to the log file.
	File string

	// Logging level, variants levels:
	//  - debug - defines debug log level
	//  - info - defines info log level
	//  - warn - defines warn log level
	//  - error - defines error log level
	//  - fatal - defines fatal log level
	//  - panic - defines panic log level
	//  - no - defines an absent log level
	//  - disabled - disables the logger
	//  - trace - defines trace log level.
	Level string
}

// Database settings.
type Database struct {
	// DSN string for database connection.
	DSN string

	// Type database, variants:
	// 	- memory
	// 	- db
	Type string
}

// RabbitMQ settings.
type RMQ struct {
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

	// Reconnection settings
	ReConnect struct {
		MaxElapsedTime  string
		InitialInterval string
		Multiplier      float64
		MaxInterval     string
	}
}

// Init is using to initialize the current config instance.
func Init(file string, c Configer) (Configer, error) {
	// read in environment variables that match
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetConfigFile(file)

	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "open config file failed")
	}

	if err := viper.Unmarshal(c); err != nil {
		return nil, errors.Wrap(err, "unmarshal config file failed")
	}

	return c, nil
}
