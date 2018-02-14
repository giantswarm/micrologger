package micrologger

import (
	"context"

	"github.com/giantswarm/microerror"
)

const (
	WildCardActivation = "*"
)

const (
	levelDebug levelID = 1 << iota
	levelInfo
	levelWarn
	levelError
)

var (
	levelMapping = map[string]levelID{
		"debug": levelDebug,
		"info":  levelInfo,
		"warn":  levelWarn,
		"error": levelError,
	}
)

type levelID byte

type ActivationLoggerConfig struct {
	Underlying Logger

	Activations map[string]string
}

type activationLogger struct {
	underlying Logger

	activations map[string]string
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

func containsKey(keyVals []interface{}, activation string) bool {
	if activation == WildCardActivation {
		return true
	}

	for i := 0; i < len(keyVals); i += 2 {
		s, ok := keyVals[i].(string)
		if ok && s == activation {
			return true
		}
	}

	return false
}

func containsVal(keyVals []interface{}, activation string) bool {
	if activation == WildCardActivation {
		return true
	}

	for i := 1; i < len(keyVals); i += 2 {
		s, ok := keyVals[i].(string)
		if ok && s == activation {
			return true
		}
	}

	return false
}

func isLevelAllowed(keyVals []interface{}, activation string) bool {
	activationLevel, ok := levelMapping[activation]
	if !ok {
		return false
	}

	for i := 0; i < len(keyVals); i += 2 {
		s, ok := keyVals[i].(string)
		if !ok {
			continue
		}
		keyValsLevel, ok := levelMapping[s]
		if !ok {
			continue
		}

		return activationLevel >= keyValsLevel
	}

	return false
}

func shouldActivate(activations map[string]string, keyVals []interface{}) (bool, error) {
	var activationCount int

	for aKey, aVal := range activations {
		if containsKey(keyVals, aKey) && containsVal(keyVals, aVal) {
			activationCount++
			continue
		}

		if isLevelAllowed(keyVals, aKey) || isLevelAllowed(keyVals, aVal) {
			activationCount++
			continue
		}
	}

	if len(activations) != 0 && len(activations) == activationCount {
		return true, nil
	}

	return false, nil
}
