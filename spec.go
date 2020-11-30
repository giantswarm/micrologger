package micrologger

import "context"

// Logger is a simple interface describing services that emit messages to
// gather certain runtime information.
type Logger interface {
	// Debug takes a format string and parameters and writes them in debug level.
	Debug(format string, params ...interface{})
	// Info takes a format string and parameters and writes them in info level.
	Info(format string, params ...interface{})
	// Warning takes a format string and parameters and writes them in warning level.
	Warning(format string, params ...interface{})
	// Error takes a format string and parameters and writes them in error level.
	Error(format string, params ...interface{})
	// Log takes a sequence of alternating key/value pairs which are used
	// to create the log message structure.
	Log(keyVals ...interface{})
	// LogCtx is the same as Log but additionally taking a context which
	// may contain additional key-value pairs that are added to the log
	// issuance, if any.
	LogCtx(ctx context.Context, keyVals ...interface{})
	// With returns a new contextual logger with keyVals appended to those
	// passed to calls to Log. If logger is also a contextual logger
	// created by With, keyVals is appended to the existing context.
	With(keyVals ...interface{}) Logger
	// WithCtx returns a context specific logger for given parameter.
	WithCtx(ctx context.Context) Logger
}
