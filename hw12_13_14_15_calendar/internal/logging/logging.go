package logging

import (
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Log struct {
	filepath, level string
	file            *os.File
}

func New(filepath, level string) *Log {
	return &Log{
		filepath: filepath,
		level:    level,
	}
}

// Init is using to initialize the Zerolog globally.
func (l *Log) Init() error {
	var err error

	l.file, err = os.OpenFile(l.filepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o666)
	if err != nil {
		return errors.Wrapf(err, "error creating file %q", l.filepath)
	}

	// Pretty logging to console
	consoleWriter := zerolog.ConsoleWriter{
		Out: os.Stderr,
	}

	// Merging log writers Zerolog output and file output
	multi := zerolog.MultiLevelWriter(consoleWriter, l.file)
	log.Logger = zerolog.New(multi).With().Timestamp().Caller().Logger()

	// Set log level. Default level in Zerolog is debug
	switch strings.ToLower(l.level) {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	case "no":
		zerolog.SetGlobalLevel(zerolog.NoLevel)
	case "disabled":
		zerolog.SetGlobalLevel(zerolog.Disabled)
	case "trace":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	}

	return nil
}

// Closing write to file.
func (l *Log) Close() error {
	return l.file.Close()
}
