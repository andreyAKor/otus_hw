package configs

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Config struct {
	// Logging settings.
	Logging struct {
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

	// Database settings
	Database struct {
		// DSN string for database connection.
		DSN string

		// Type database, variants:
		// 	- memory
		// 	- db
		Type string
	}

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

// Init is using to initialize the configs.
func Init(file string) (Config, error) {
	// read in environment variables that match
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetConfigFile(file)

	var c Config

	if err := viper.ReadInConfig(); err != nil {
		return c, errors.Wrap(err, "open config file failed")
	}

	if err := viper.Unmarshal(&c); err != nil {
		return c, errors.Wrap(err, "unmarshal config file failed")
	}

	return c, nil
}
