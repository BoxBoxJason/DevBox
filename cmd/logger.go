package main

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// SetupZapLogger sets up the Zap logger with the specified verbosity level and optional file output.
// It configures the logger to use ISO8601 time format and capitalizes the log levels.
// The logger is set to production mode by default, but can be configured for debug mode if verbose is true.
// If filename is not empty, the logger's output is redirected to the specified file.
func SetupZapLogger(verbose bool, filename string) {
	// Set up the logger configuration
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	if verbose {
		config.Level.SetLevel(zapcore.DebugLevel)
	} else {
		config.Level.SetLevel(zapcore.InfoLevel)
	}

	// If a filename is specified, update the output path.
	if filename != "" {
		err := os.MkdirAll(filepath.Dir(filename), 0700)
		if err != nil {
			zap.L().Fatal("Failed to create log directory: " + err.Error())
		}
		config.OutputPaths = []string{filename, "stderr"}
	}

	// Create the logger
	logger, err := config.Build()
	if err != nil {
		zap.L().Fatal("Failed to create logger: " + err.Error())
	}

	// Set the global logger
	zap.ReplaceGlobals(logger)
}
