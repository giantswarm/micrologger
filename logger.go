// Package logger implements a logging interface used to log messages.
package micrologger

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"

	kitlog "github.com/go-kit/kit/log"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger/loggermeta"
)

type Config struct {
	Caller             kitlog.Valuer
	IOWriter           io.Writer
	TimestampFormatter kitlog.Valuer
}

type MicroLogger struct {
	logger kitlog.Logger
}

type CtxMicroLogger struct {
	ctx        context.Context
	underlying Logger
}

func New(config Config) (*MicroLogger, error) {
	if config.Caller == nil {
		config.Caller = DefaultCaller
	}
	if config.TimestampFormatter == nil {
		config.TimestampFormatter = DefaultTimestampFormatter
	}
	if config.IOWriter == nil {
		config.IOWriter = DefaultIOWriter
	}

	kitLogger := kitlog.NewJSONLogger(kitlog.NewSyncWriter(config.IOWriter))
	kitLogger = kitlog.With(
		kitLogger,
		"caller", config.Caller,
		"time", config.TimestampFormatter,
	)

	l := &MicroLogger{
		logger: kitLogger,
	}

	return l, nil
}

func (l *MicroLogger) Debug(format string, params ...interface{}) {
	l.Log("level", "debug", "message", fmt.Sprintf(format, params...))
}

func (l *MicroLogger) Info(format string, params ...interface{}) {
	l.Log("level", "info", "message", fmt.Sprintf(format, params...))
}

func (l *MicroLogger) Warning(format string, params ...interface{}) {
	l.Log("level", "warning", "message", fmt.Sprintf(format, params...))
}

func (l *MicroLogger) Error(err error, format string, params ...interface{}) {
	l.Log("level", "error", "message", fmt.Sprintf(format, params...), "stack", microerror.JSON(err))
}

func (l *MicroLogger) Log(keyVals ...interface{}) {
	keyVals = l.processStack(keyVals)
	err := l.logger.Log(keyVals...)
	if err != nil {
		log.Printf("failed to log, reason: %#q", err.Error())
	}
}

func (l *MicroLogger) LogCtx(ctx context.Context, keyVals ...interface{}) {
	keyVals = l.processStack(keyVals)
	meta, ok := loggermeta.FromContext(ctx)
	if !ok {
		err := l.logger.Log(keyVals...)
		if err != nil {
			log.Printf("failed to log, reason: %#q", err.Error())
		}
		return
	}

	var newKeyVals []interface{}
	{
		newKeyVals = append(newKeyVals, keyVals...)

		for k, v := range meta.KeyVals {
			newKeyVals = append(newKeyVals, k)
			newKeyVals = append(newKeyVals, v)
		}
	}

	err := l.logger.Log(newKeyVals...)
	if err != nil {
		log.Printf("failed to log, reason: %#q", err.Error())
	}
}

func (l *MicroLogger) With(keyVals ...interface{}) Logger {
	keyVals = l.processStack(keyVals)
	return &MicroLogger{
		logger: kitlog.With(l.logger, keyVals...),
	}
}

func (l *MicroLogger) WithCtx(ctx context.Context) Logger {
	return &CtxMicroLogger{
		ctx:        ctx,
		underlying: l,
	}
}

func (l *MicroLogger) processStack(keyVals []interface{}) []interface{} {
	for i := 1; i < len(keyVals); i += 2 {
		k := keyVals[i-1]
		v := keyVals[i]

		// If this is not the "stack" key try on next iteration.
		if k != "stack" {
			continue
		}

		// Try to get bytes of the data for the "stack" key. Return
		// what is given otherwise.
		var bytes []byte
		switch data := v.(type) {
		case string:
			bytes = []byte(data)
		case []byte:
			bytes = data
		default:
			return keyVals
		}

		// If the found value isn't a JSON return.
		var m map[string]interface{}
		err := json.Unmarshal(bytes, &m)
		if err != nil {
			return keyVals
		}

		// If the found value is a JSON then make a copy of keyVals to
		// not mutate the original one and store the value as a map to
		// be rendered as a JSON object. Then return it.
		keyValsCopy := append([]interface{}{}, keyVals...)
		keyValsCopy[i] = m

		return keyValsCopy
	}

	return keyVals
}

func (l *CtxMicroLogger) Debug(format string, params ...interface{}) {
	l.LogCtx(l.ctx, "level", "debug", "message", fmt.Sprintf(format, params...))
}

func (l *CtxMicroLogger) Info(format string, params ...interface{}) {
	l.LogCtx(l.ctx, "level", "info", "message", fmt.Sprintf(format, params...))
}

func (l *CtxMicroLogger) Warning(format string, params ...interface{}) {
	l.LogCtx(l.ctx, "level", "warning", "message", fmt.Sprintf(format, params...))
}

func (l *CtxMicroLogger) Error(err error, format string, params ...interface{}) {
	l.LogCtx(l.ctx, "level", "error", "message", fmt.Sprintf(format, params...), "stack", microerror.JSON(err))
}

func (l *CtxMicroLogger) Log(keyVals ...interface{}) {
	l.underlying.Log(keyVals...)
}

func (l *CtxMicroLogger) LogCtx(ctx context.Context, keyVals ...interface{}) {
	l.underlying.LogCtx(ctx, keyVals...)
}

func (l *CtxMicroLogger) With(keyVals ...interface{}) Logger {
	return l.underlying.With(keyVals...)
}

func (l *CtxMicroLogger) WithCtx(ctx context.Context) Logger {
	return &CtxMicroLogger{
		ctx:        ctx,
		underlying: l,
	}
}
