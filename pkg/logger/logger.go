package logger

import (
	"fmt"
	"time"

	"github.com/rs/zerolog"
)

func init() {
	zerolog.TimestampFieldName = "ts"
	zerolog.MessageFieldName = "msg"
	zerolog.TimestampFunc = func() time.Time {
		return time.Now().UTC()
	}
}

type RotatorCloser interface {
	Rotate() error
	Close() error
}

func NewReloadableLogger(filename, level string) (*zerolog.Logger, RotatorCloser, error) {
	switch level {
	case "disabled":
		zerolog.SetGlobalLevel(zerolog.Disabled)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	default:
		return nil, nil, fmt.Errorf("invalid log level: %s", level)
	}

	w, err := newLogWriter(filename)
	if err != nil {
		return nil, nil, err
	}

	l := zerolog.New(w).With().Timestamp().Logger()

	return &l, w, nil
}
