// Package logger implements a logging interface used to log messages.
package micrologger

import (
	"context"
	"io"

	kitlog "github.com/go-kit/kit/log"

	"github.com/giantswarm/micrologger/loggermeta"
)

type HumanConfig struct {
	Caller             kitlog.Valuer
	IOWriter           io.Writer
	TimestampFormatter kitlog.Valuer

	Verbose bool
}

type HumanLogger struct {
	logger kitlog.Logger
}

func NewHumanLogger(config Config) (*HumanLogger, error) {
	if config.Caller == nil {
		config.Caller = DefaultCaller
	}
	if config.TimestampFormatter == nil {
		config.TimestampFormatter = DefaultTimestampFormatter
	}
	if config.IOWriter == nil {
		config.IOWriter = DefaultIOWriter
	}

	kitLogger := kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(config.IOWriter))

	kitLogger = kitlog.With(
		kitLogger,
		"caller", config.Caller,
		"time", config.TimestampFormatter,
	)

	l := &HumanLogger{
		logger: kitLogger,
	}

	return l, nil
}

func (l *HumanLogger) Log(keyVals ...interface{}) error {
	var output []interface{}
	for i := 0; i < len(keyVals); i += 2 {
		k := keyVals[i]
		v := keyVals[i+1]
		switch k {
		case "message":
			//output = append(output, k)
			output = append(output, v)
		}
	}
	return l.logger.Log(output)
}

func (l *HumanLogger) LogCtx(ctx context.Context, keyVals ...interface{}) error {
	meta, ok := loggermeta.FromContext(ctx)
	if !ok {
		return l.logger.Log(keyVals...)
	}

	var newKeyVals []interface{}
	{
		newKeyVals = append(newKeyVals, keyVals...)

		for k, v := range meta.KeyVals {
			newKeyVals = append(newKeyVals, k)
			newKeyVals = append(newKeyVals, v)
		}
	}

	return l.logger.Log(newKeyVals...)
}

func (l *HumanLogger) With(keyVals ...interface{}) Logger {
	return &HumanLogger{
		logger: kitlog.With(l.logger, keyVals...),
	}
}
