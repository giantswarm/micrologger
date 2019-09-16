// Package logger implements a logging interface used to log messages.
package micrologger

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/fatih/color"

	"github.com/giantswarm/micrologger/loggermeta"
)

type Valuer func() interface{}

type Config struct {
	Caller             Valuer
	IOWriter           io.Writer
	TimestampFormatter Valuer
	Human              bool
}

type MicroLogger struct {
	human   bool
	log     *log.Logger // store as pointer as logger contains a mutex
	keyVals map[interface{}]interface{}
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

	log := log.New(config.IOWriter, "", 0)

	logger := &MicroLogger{
		human: config.Human,
		keyVals: map[interface{}]interface{}{
			"caller": config.Caller,
			"time":   config.TimestampFormatter,
			"level":  "info",
		},
		log: log,
	}

	return logger, nil
}

func ToString(v interface{}) string {
	valuer, ok := v.(Valuer)
	if ok {
		v = valuer()
	}
	return fmt.Sprintf("%v", v)
}

func (l *MicroLogger) Log(keyVals ...interface{}) error {
	combined := map[string]string{}
	for k, v := range l.keyVals {
		k := ToString(k)
		v := ToString(v)
		combined[k] = v
	}
	for i := 0; i < len(keyVals); i += 2 {
		k := ToString(keyVals[i])
		v := ToString(keyVals[i+1])
		combined[k] = v
	}
	var encoded string
	if l.human {
		level := strings.ToUpper(combined["level"])
		switch level {
		case "ERROR":
			level = color.RedString(level)
			break
		case "WARN":
			level = color.YellowString(level)
			break
		}
		encoded = fmt.Sprintf("%s %s %s", color.GreenString(combined["time"]), level, combined["message"])
	} else {
		marshalled, err := json.Marshal(combined)
		if err != nil {
			return err
		}
		encoded = string(marshalled)
	}
	l.log.Println(encoded)
	return nil
}

func (l *MicroLogger) LogCtx(ctx context.Context, keyVals ...interface{}) error {
	meta, ok := loggermeta.FromContext(ctx)
	if !ok {
		return l.Log(keyVals...)
	}

	combined := keyVals
	for k, v := range meta.KeyVals {
		combined = append(combined, k)
		combined = append(combined, v)
	}

	return l.Log(combined...)
}

func (l *MicroLogger) With(keyVals ...interface{}) Logger {
	combined := map[interface{}]interface{}{}
	for k, v := range l.keyVals {
		combined[k] = v
	}
	for i := 0; i < len(keyVals); i += 2 {
		k := keyVals[i]
		v := keyVals[i+1]
		combined[k] = v
	}
	return &MicroLogger{
		keyVals: combined,
		log:     l.log,
	}
}
