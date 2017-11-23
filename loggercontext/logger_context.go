// Package loggercontext stores and accesses the container struct in
// context.Context.
package loggercontext

import (
	"context"
)

// key is an unexported type for keys defined in this package. This prevents
// collisions with keys defined in other packages.
type key string

// loggerKey is the key for logger struct values in context.Context. Clients use
// loggercontext.NewContext and loggercontext.FromContext instead of using this
// key directly.
var loggerKey key = "logger"

// LoggerContext is a communication structure used to transport information in order
// for a micro logger to use it when issuing logs.
type LoggerContext struct {
	// KeyVals is a mapping of key-value pairs a micro logger adds to the log
	// message issuance.
	KeyVals map[string]string
}

func New() *LoggerContext {
	return &LoggerContext{
		KeyVals: map[string]string{},
	}
}

// NewContext returns a new context.Context that carries value v.
func NewContext(ctx context.Context, v *LoggerContext) context.Context {
	if v == nil {
		return ctx
	}

	return context.WithValue(ctx, loggerKey, v)
}

// FromContext returns the logger struct, if any.
func FromContext(ctx context.Context) (*LoggerContext, bool) {
	v, ok := ctx.Value(loggerKey).(*LoggerContext)
	return v, ok
}
