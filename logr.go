package micrologger

import (
	"context"
	"strings"

	kitlog "github.com/go-kit/log"
	"github.com/go-logr/logr"
	"github.com/go-stack/stack"
)

type LogrSink struct {
	*MicroLogger
}

func (l *LogrSink) Init(info logr.RuntimeInfo) {
	l.info = info
}

func (l *LogrSink) Enabled(level int) bool {
	return l.verbosity >= level
}

func (l *LogrSink) Info(level int, msg string, keysAndValues ...interface{}) {
	if l.verbosity < level {
		return
	}
	keysAndValues = append(keysAndValues, l.getValues("debug")...)
	l.With(keysAndValues...).Log("message", msg)
}

func (l *LogrSink) Error(err error, msg string, keysAndValues ...interface{}) {
	if l.verbosity < 1 {
		return
	}
	keysAndValues = append(keysAndValues, l.getValues("error")...)
	l.With(keysAndValues...).Errorf(context.Background(), err, msg)
}

func (l *LogrSink) WithValues(keysAndValues ...interface{}) logr.LogSink {
	loggerCopy := l.deepCopy()
	loggerCopy.logger = kitlog.With(loggerCopy.logger, processStack(keysAndValues)...)
	return loggerCopy.AsSink(loggerCopy.verbosity)
}

func (l *LogrSink) WithName(name string) logr.LogSink {
	loggerCopy := l.deepCopy()
	loggerCopy.names = append(loggerCopy.names[:], name)
	return loggerCopy.AsSink(l.verbosity)
}

func (l *LogrSink) getValues(level string) []interface{} {
	return []interface{}{
		"level",
		level,
		"name",
		strings.Join(l.names, "."),
		"caller",
		stack.Caller(l.info.CallDepth),
	}
}
