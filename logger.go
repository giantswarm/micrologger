// Package logger implements a logging interface used to log messages.
package micrologger

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/giantswarm/microerror"
	kitlog "github.com/go-kit/kit/log"

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

func (l *MicroLogger) Debugf(ctx context.Context, format string, params ...interface{}) {
	kvs := []interface{}{
		"level", "debug",
		"message", fmt.Sprintf(format, params...),
	}

	l.log(keyValsWithMeta(ctx, kvs))
}

func (l *MicroLogger) Errorf(ctx context.Context, err error, format string, params ...interface{}) {
	var kvs []interface{}
	if err != nil {
		kvs = []interface{}{
			"level", "error",
			"message", fmt.Sprintf(format, params...),
			"stack", microerror.JSON(err),
		}
	} else {
		kvs = []interface{}{
			"level", "error",
			"message", fmt.Sprintf(format, params...),
		}
	}

	l.log(keyValsWithMeta(ctx, kvs))
}

func (l *MicroLogger) Log(keyVals ...interface{}) {
	l.log(processStack(keyVals))
}

func (l *MicroLogger) LogCtx(ctx context.Context, keyVals ...interface{}) {
	l.log(keyValsWithMeta(ctx, keyVals))
}

func (l *MicroLogger) With(keyVals ...interface{}) Logger {
	keyVals = processStack(keyVals)
	return &MicroLogger{
		logger: kitlog.With(l.logger, keyVals...),
	}
}

func (l *MicroLogger) log(keyVals []interface{}) {
	err := l.logger.Log(keyVals...)
	if err != nil {
		log.Printf("failed to log with error: %#q, keyVals = %v", err.Error(), keyVals)
	}
}

func keyValsWithMeta(ctx context.Context, keyVals []interface{}) []interface{} {
	keyVals = processStack(keyVals)
	meta, ok := loggermeta.FromContext(ctx)
	if !ok {
		return keyVals
	}

	var kvs []interface{}
	{
		kvs = append(kvs, keyVals...)

		for k, v := range meta.KeyVals {
			kvs = append(kvs, k)
			kvs = append(kvs, v)
		}
	}

	return kvs
}

func processStack(keyVals []interface{}) []interface{} {
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
