package main

import (
	"fmt"
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

	// Custom time encoder to include only date, time (hour, minute, second, millisecond)
	config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05(000)")

	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	// Disable automatic caller/file:line and stacktraces
	config.DisableCaller = true
	config.DisableStacktrace = true

	if verbose {
		config.Level.SetLevel(zapcore.DebugLevel)
	} else {
		config.Level.SetLevel(zapcore.InfoLevel)
	}

	// If a filename is specified, update the output path.
	if filename != "" {
		if err := os.MkdirAll(filepath.Dir(filename), 0700); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create log directory: %v\n", err)
			os.Exit(1)
		}
		config.OutputPaths = []string{filename, "stderr"}
	}

	// Create the logger
	logger, err := config.Build()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create logger: %v\n", err)
		os.Exit(1)
	}

	// Set the global logger
	zap.ReplaceGlobals(logger)
}
