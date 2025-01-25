package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
)

var log zerolog.Logger

func Init(logDir string) error {
	// Create logs directory if not exists
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %v", err)
	}

	// Create log file with timestamp
	logFile := filepath.Join(logDir, fmt.Sprintf("neon-%s.log", time.Now().Format("2006-01-02")))
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %v", err)
	}

	// Configure zerolog
	zerolog.TimeFieldFormat = time.RFC3339
	log = zerolog.New(file).With().Timestamp().Logger()

	log.Info().Msg("Logger initialized")
	return nil
}

// Helper functions for different log levels
func Info(msg string) {
	log.Info().Msg(msg)
}

func Infof(format string, v ...interface{}) {
	log.Info().Msgf(format, v...)
}

func Error(err error, msg string) {
	log.Error().Err(err).Msg(msg)
}

func Errorf(err error, format string, v ...interface{}) {
	log.Error().Err(err).Msgf(format, v...)
}

func Debug(msg string) {
	log.Debug().Msg(msg)
}

func Debugf(format string, v ...interface{}) {
	log.Debug().Msgf(format, v...)
}

func Warn(msg string) {
	log.Warn().Msg(msg)
}

func Warnf(format string, v ...interface{}) {
	log.Warn().Msgf(format, v...)
}

// WithField adds a field to the log entry
func WithField(key string, value interface{}) zerolog.Logger {
	return log.With().Interface(key, value).Logger()
}
