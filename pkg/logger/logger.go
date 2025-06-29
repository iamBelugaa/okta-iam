// Package logger provides a centralized and configurable logging utility
// based on the `go.uber.org/zap` library. It offers a structured logging
// approach suitable for production environments, allowing for easy integration
// and consistent log formatting across different services and applications.
// The package simplifies the setup of a robust logger with predefined
// configurations for output format, timestamps, and initial contextual fields.
package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Creates and configures a new Zap SugaredLogger.
// It sets up a production-ready logger with JSON encoding, ISO8601 timestamps,
// and includes service name and process ID as initial fields.
func New(service string, outputPaths ...string) *zap.SugaredLogger {
	encoderCfg := zap.NewProductionEncoderConfig()

	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	// Initialize the Zap configuration. This struct holds all the settings
	// for building the logger.
	config := zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:       false,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling:          nil,
		Encoding:          "json",
		EncoderConfig:     encoderCfg,
		OutputPaths:       []string{"stderr"},
		ErrorOutputPaths:  []string{"stderr"},
		InitialFields:     map[string]any{"service": service, "pid": os.Getpid()},
	}

	// Override the default output paths if custom 'outputPaths' are provided.
	if len(outputPaths) != 0 {
		config.OutputPaths = outputPaths
	}

	return zap.Must(config.Build()).Sugar()
}
