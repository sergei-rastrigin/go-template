package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

type contextKeyType int

var (
	KeyLog   = contextKeyType(1)
	KeyTrace = contextKeyType(2)
)

var defaultLogger = zerolog.New(os.Stderr).With().Timestamp().Logger().Output(zerolog.ConsoleWriter{
	Out:        os.Stderr,
	TimeFormat: time.RFC3339Nano,
})

type LogConfig struct {
	Level  string `envconfig:"LOG_LEVEL" default:"info"`
	Output string `envconfig:"LOG_OUTPUT" default:"console"`
}

func NewLogger(cfg LogConfig, serviceName string) (zerolog.Logger, error) {
	level, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		return zerolog.Logger{}, err
	}

	logger := defaultLogger.With().Str("service", serviceName).Logger().Level(level)

	var writer io.Writer

	if cfg.Output == "console" {
		writer = zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: time.RFC3339Nano,
		}
		logger = logger.Output(writer)
	}

	return logger, nil
}
