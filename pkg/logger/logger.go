package logger

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func ConfigureLogger() {
	// Configure zerolog
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logLevel := zerolog.InfoLevel // Default level

	if lvl, found := os.LookupEnv("LOG_LEVEL"); found {
		level, err := zerolog.ParseLevel(lvl)
		if err == nil {
			logLevel = level
		}
	}
	zerolog.SetGlobalLevel(logLevel)

	// Configure console output
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stderr}
	log.Logger = zerolog.New(consoleWriter).With().Timestamp().Logger()
}
