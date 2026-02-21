// Package telemetry provides structured logging for the NetSentry application.
package telemetry

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

// Logger is the application-wide structured logger.
type Logger struct {
	zl zerolog.Logger
}

// LogOptions configures the logger.
type LogOptions struct {
	// Level is the minimum log level: "debug", "info", "warn", "error".
	Level string
	// JSON enables JSON-formatted log output. Defaults to human-readable console output.
	JSON bool
	// Output is the log destination. Defaults to os.Stderr.
	Output io.Writer
}

// NewLogger constructs a Logger with the given options.
func NewLogger(opts LogOptions) *Logger {
	out := opts.Output
	if out == nil {
		out = os.Stderr
	}

	var zl zerolog.Logger
	if opts.JSON {
		zl = zerolog.New(out).With().Timestamp().Logger()
	} else {
		cw := zerolog.ConsoleWriter{
			Out:        out,
			TimeFormat: time.RFC3339,
		}
		zl = zerolog.New(cw).With().Timestamp().Logger()
	}

	level := zerolog.InfoLevel
	switch opts.Level {
	case "debug":
		level = zerolog.DebugLevel
	case "warn":
		level = zerolog.WarnLevel
	case "error":
		level = zerolog.ErrorLevel
	case "trace":
		level = zerolog.TraceLevel
	}
	zl = zl.Level(level)

	return &Logger{zl: zl}
}

// With returns a new Logger with the given key-value fields attached.
func (l *Logger) With(key, value string) *Logger {
	return &Logger{zl: l.zl.With().Str(key, value).Logger()}
}

// Debug logs a debug message.
func (l *Logger) Debug(msg string, fields ...interface{}) {
	l.event(l.zl.Debug(), msg, fields...)
}

// Info logs an informational message.
func (l *Logger) Info(msg string, fields ...interface{}) {
	l.event(l.zl.Info(), msg, fields...)
}

// Warn logs a warning message.
func (l *Logger) Warn(msg string, fields ...interface{}) {
	l.event(l.zl.Warn(), msg, fields...)
}

// Error logs an error message.
func (l *Logger) Error(msg string, err error, fields ...interface{}) {
	ev := l.zl.Error()
	if err != nil {
		ev = ev.Err(err)
	}
	l.event(ev, msg, fields...)
}

// Fatal logs a fatal message and terminates the process.
func (l *Logger) Fatal(msg string, err error) {
	l.zl.Fatal().Err(err).Msg(msg)
}

// event applies key-value pairs to a zerolog.Event and sends the message.
// Fields must be alternating string key-value pairs.
func (l *Logger) event(ev *zerolog.Event, msg string, fields ...interface{}) {
	for i := 0; i+1 < len(fields); i += 2 {
		k, _ := fields[i].(string)
		ev = ev.Interface(k, fields[i+1])
	}
	ev.Msg(msg)
}
