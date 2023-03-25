package micrologger

import "context"

// Logger is a simple interface describing services that emit messages to
// gather certain runtime information.
type Logger interface {
	// Debug writes the given message in debug level.
	Debug(keyVals ...interface{})
	// DebugCtx is the same as Debug but additionally taking a context which
	// may contain additional key-value pairs that are added to the log
	// issuance, if any.
	DebugCtx(ctx context.Context, keyVals ...interface{})
	// Debugf takes a format string and parameters and writes them in debug level.
	Debugf(ctx context.Context, format string, params ...interface{})
	// Error takes an error and a message and writes them in error level. The
	// error stack trace is written as "stack" value log entry.
	Error(err error, keyVals ...interface{})
	// ErrorCtx is the same as Error but additionally taking a context which
	// may contain additional key-value pairs that are added to the log
	// issuance, if any.
	ErrorCtx(ctx context.Context, err error, keyVals ...interface{})
	// Errorf takes an error, a format string and parameters and writes them in
	// error level. The error stack trace is written as "stack" value log
	// entry.
	Errorf(ctx context.Context, err error, format string, params ...interface{})
	// Info writes the given message in info level.
	Info(keyVals ...interface{})
	// InfoCtx is the same as Info but additionally taking a context which
	// may contain additional key-value pairs that are added to the log
	// issuance, if any.
	InfoCtx(ctx context.Context, keyVals ...interface{})
	// Infof takes a format string and parameters and writes them in info level.
	Infof(ctx context.Context, format string, params ...interface{})
	// Log takes a sequence of alternating key/value pairs which are used
	// to create the log message structure.
	Log(keyVals ...interface{})
	// LogCtx is the same as Log but additionally taking a context which
	// may contain additional key-value pairs that are added to the log
	// issuance, if any.
	LogCtx(ctx context.Context, keyVals ...interface{})
	// SetLevel configures log level.
	SetLevel(lvl string)
	// Warning writes the given message in warning level.
	Warning(keyVals ...interface{})
	// WarningCtx is the same as Warning but additionally taking a context which
	// may contain additional key-value pairs that are added to the log
	// issuance, if any.
	WarningCtx(ctx context.Context, keyVals ...interface{})
	// Warningf takes a format string and parameters and writes them in warning level.
	Warningf(ctx context.Context, format string, params ...interface{})
	// With returns a new contextual logger with keyVals appended to those
	// passed to calls to Log. If logger is also a contextual logger
	// created by With, keyVals is appended to the existing context.
	With(keyVals ...interface{}) Logger
	// WithIncreasedCallerDepth is useful when wrapping with another
	// interface to pass it as dependency to a library outside Giant Swarm.
	WithIncreasedCallerDepth() Logger
}
