package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Level is a log level supported by the logger.
type Level string

const (
	DebugLevel Level = "debug"
	InfoLevel  Level = "info"
	WarnLevel  Level = "warn"
	ErrorLevel Level = "error"
)

// config holds logger configuration.
type config struct {
	level       Level
	development bool
}

// Option configures a logger.
type Option func(*config)

// WithLevel sets the log level.
func WithLevel(l Level) Option {
	return func(c *config) {
		c.level = l
	}
}

// WithDevelopment enables or disables development mode.
func WithDevelopment(dev bool) Option {
	return func(c *config) {
		c.development = dev
	}
}

// New creates a zap.Logger based on the provided options.
// It writes structured logs to stderr and returns an error for unsupported levels.
// Defaults to InfoLevel if no level is provided.
func New(opts ...Option) (*zap.Logger, error) {
	cfg := config{
		level: InfoLevel,
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	level, err := ParseLevel(cfg.level)
	if err != nil {
		return nil, fmt.Errorf("parse log level: %w", err)
	}

	var zapCfg zap.Config

	switch cfg.development {
	case true:
		zapCfg = zap.NewDevelopmentConfig()
		zapCfg.Level = zap.NewAtomicLevelAt(level)
		zapCfg.OutputPaths = []string{"stderr"}
		zapCfg.ErrorOutputPaths = []string{"stderr"}
	default:
		zapCfg = zap.NewProductionConfig()
		zapCfg.Level = zap.NewAtomicLevelAt(level)
		zapCfg.OutputPaths = []string{"stderr"}
		zapCfg.ErrorOutputPaths = []string{"stderr"}
	}

	logger, err := zapCfg.Build()
	if err != nil {
		return nil, fmt.Errorf("build logger: %w", err)
	}

	return logger, nil
}

func ParseLevel(level Level) (zapcore.Level, error) {
	switch level {
	case DebugLevel:
		return zapcore.DebugLevel, nil
	case InfoLevel:
		return zapcore.InfoLevel, nil
	case WarnLevel:
		return zapcore.WarnLevel, nil
	case ErrorLevel:
		return zapcore.ErrorLevel, nil
	default:
		return zapcore.InfoLevel, fmt.Errorf("unsupported log level: %s", level)
	}
}
