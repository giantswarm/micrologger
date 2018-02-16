package micrologger

import (
	"context"

	"github.com/giantswarm/microerror"
)

const (
	KeyLevel     = "level"
	KeyVerbosity = "verbosity"
)

const (
	levelDebug levelID = 1 << iota
	levelInfo
	levelWarning
	levelError
)

var (
	levelMapping = map[string]levelID{
		"debug":   levelDebug,
		"info":    levelInfo,
		"warning": levelWarning,
		"error":   levelError,
	}
)

type levelID byte

type ActivationLoggerConfig struct {
	Underlying Logger

	Activations map[string]interface{}
}

type activationLogger struct {
	underlying Logger

	activations map[string]interface{}
}

func NewActivation(config ActivationLoggerConfig) (Logger, error) {
	if config.Underlying == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Underlying must not be empty", config)
	}

	l := &activationLogger{
		underlying: config.Underlying,

		activations: config.Activations,
	}

	return l, nil
}

func (l *activationLogger) Log(keyVals ...interface{}) error {
	activated, err := shouldActivate(l.activations, keyVals)
	if err != nil {
		return microerror.Mask(err)
	}

	if activated {
		return l.underlying.Log(keyVals...)
	}

	return nil
}

func (l *activationLogger) LogCtx(ctx context.Context, keyVals ...interface{}) error {
	activated, err := shouldActivate(l.activations, keyVals)
	if err != nil {
		return microerror.Mask(err)
	}

	if activated {
		return l.underlying.LogCtx(ctx, keyVals...)
	}

	return nil
}

func (l *activationLogger) With(keyVals ...interface{}) Logger {
	return l.underlying.With(keyVals...)
}

func containsKey(keyVals []interface{}, aKey string) bool {
	for i := 0; i < len(keyVals); i += 2 {
		s, ok := keyVals[i].(string)
		if ok && s == aKey {
			return true
		}
	}

	return false
}

func containsVal(keyVals []interface{}, aVal interface{}) bool {
	for i := 1; i < len(keyVals); i += 2 {
		if keyVals[i] == aVal {
			return true
		}
	}

	return false
}

func isLevelAllowed(keyVals []interface{}, aVal interface{}) bool {
	s, ok := aVal.(string)
	if !ok {
		return false
	}
	activationLevel, ok := levelMapping[s]
	if !ok {
		return false
	}

	for i := 0; i < len(keyVals); i += 2 {
		k, ok := keyVals[i].(string)
		if !ok {
			continue
		}
		if k != KeyLevel {
			continue
		}
		v, ok := keyVals[i+1].(string)
		if !ok {
			continue
		}
		keyValsLevel, ok := levelMapping[v]
		if !ok {
			continue
		}

		return activationLevel >= keyValsLevel
	}

	return false
}

func isVerbosityAllowed(keyVals []interface{}, aVal interface{}) bool {
	activationVerbosity, ok := aVal.(int)
	if !ok {
		return false
	}

	for i := 0; i < len(keyVals); i += 2 {
		k, ok := keyVals[i].(string)
		if !ok {
			continue
		}
		if k != KeyVerbosity {
			continue
		}
		keyValsVerbosity, ok := keyVals[i+1].(int)
		if !ok {
			continue
		}

		return activationVerbosity >= keyValsVerbosity
	}

	return false
}

func shouldActivate(activations map[string]interface{}, keyVals []interface{}) (bool, error) {
	var activationCount int

	for aKey, aVal := range activations {
		if containsKey(keyVals, aKey) && containsVal(keyVals, aVal) {
			activationCount++
			continue
		}
		if aKey == KeyLevel && isLevelAllowed(keyVals, aVal) {
			activationCount++
			continue
		}
		if aKey == KeyVerbosity && isVerbosityAllowed(keyVals, aVal) {
			activationCount++
			continue
		}
	}

	if len(activations) != 0 && len(activations) == activationCount {
		return true, nil
	}

	return false, nil
}
