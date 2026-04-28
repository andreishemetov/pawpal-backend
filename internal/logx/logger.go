package logx

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

func New() zerolog.Logger {
	logger := zerolog.New(os.Stdout).
		With().
		Timestamp().
		Logger()

	zerolog.TimeFieldFormat = time.RFC3339

	return logger
}