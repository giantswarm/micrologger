// Package logger implements a logging interface used to log messages.
package micrologger

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/giantswarm/microerror"
	kitlog "github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/go-logr/logr"

	"github.com/giantswarm/micrologger/loggermeta"
)

type Config struct {
	Caller   kitlog.Valuer
	IOWriter io.Writer
	// OutputFormat can be either text or json (default).
	OutputFormat       string
	TimestampFormatter kitlog.Valuer
}

type MicroLogger struct {
	info      logr.RuntimeInfo
	logger    kitlog.Logger
	verbosity int
	names     []string
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

	var kitLogger kitlog.Logger
	switch strings.ToLower(config.OutputFormat) {
	case "text":
		kitLogger = kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(config.IOWriter))
	default:
		kitLogger = kitlog.NewJSONLogger(kitlog.NewSyncWriter(config.IOWriter))
	}

	//kitLogger = kitlog.WithSuffix(kitLogger, "caller", kitlog.Caller(8))

	logger := kitlog.LoggerFunc(func(keyVals ...interface{}) error {
		err := kitLogger.Log(keyVals...)
		if err != nil {
			log.Printf("failed to log with error: %#q, keyVals = %v", err.Error(), keyVals)
		}
		return nil
	})

	l := &MicroLogger{
		logger: logger,
	}

	return l, nil
}

func (l *MicroLogger) SetLevel(lvl string) {
	l.logger = level.NewFilter(l.logger, level.Allow(level.ParseDefault(lvl, level.InfoValue())))
}

func (l *MicroLogger) Debug(keyVals ...interface{}) {
	level.Debug(l.logger).Log(keyVals...)
}

func (l *MicroLogger) DebugCtx(ctx context.Context, keyVals ...interface{}) {
	l.Debug(keyValsWithMeta(ctx, keyVals)...)
}

func (l *MicroLogger) Debugf(ctx context.Context, format string, params ...interface{}) {
	l.DebugCtx(ctx, "message", fmt.Sprintf(format, params...))
}

func (l *MicroLogger) Error(err error, keyVals ...interface{}) {
	errKV := []interface{}{
		"error", microerror.JSON(err),
	}
	keyVals = append(keyVals, errKV...)

	level.Error(l.logger).Log(keyVals...)
}

func (l *MicroLogger) ErrorCtx(ctx context.Context, err error, keyVals ...interface{}) {
	l.Error(err, keyValsWithMeta(ctx, keyVals)...)
}

func (l *MicroLogger) Errorf(ctx context.Context, err error, format string, params ...interface{}) {
	l.ErrorCtx(ctx, err, "message", fmt.Sprintf(format, params...))
}

func (l *MicroLogger) Info(keyVals ...interface{}) {
	level.Warn(l.logger).Log(keyVals...)
}

func (l *MicroLogger) InfoCtx(ctx context.Context, keyVals ...interface{}) {
	l.Info(keyValsWithMeta(ctx, keyVals)...)
}

func (l *MicroLogger) Infof(ctx context.Context, format string, params ...interface{}) {
	l.InfoCtx(ctx, "message", fmt.Sprintf(format, params...))
}

func (l *MicroLogger) Log(keyVals ...interface{}) {
	l.logger.Log(processStack(keyVals)...)
}

func (l *MicroLogger) LogCtx(ctx context.Context, keyVals ...interface{}) {
	l.logger.Log(keyValsWithMeta(ctx, keyVals)...)
}

func (l *MicroLogger) Warning(keyVals ...interface{}) {
	level.Warn(l.logger).Log(keyVals...)
}

func (l *MicroLogger) WarningCtx(ctx context.Context, keyVals ...interface{}) {
	l.Warning(keyValsWithMeta(ctx, keyVals)...)
}

func (l *MicroLogger) Warningf(ctx context.Context, format string, params ...interface{}) {
	l.WarningCtx(ctx, "message", fmt.Sprintf(format, params...))
}

func (l *MicroLogger) deepCopy() *MicroLogger {
	return &MicroLogger{
		info:      l.info,
		logger:    l.logger,
		verbosity: l.verbosity,
		names:     l.names[:],
	}
}

func (l *MicroLogger) With(keyVals ...interface{}) Logger {
	loggerCopy := l.deepCopy()
	loggerCopy.logger = kitlog.With(loggerCopy.logger, keyVals...)
	return loggerCopy
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

func (l *MicroLogger) WithIncreasedCallerDepth() Logger {
	return &MicroLogger{
		logger: kitlog.With(l.logger, "caller", newCallerFunc(1)),
	}
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
