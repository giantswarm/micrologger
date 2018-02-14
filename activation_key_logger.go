package micrologger

import (
	"context"

	"github.com/giantswarm/microerror"
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

type ActivationKeyLoggerConfig struct {
	Underlying Logger

	ActivationKeys []string
}

type activationKeyLogger struct {
	underlying Logger

	activationKeys []string
}

func NewActivationKey(config ActivationKeyLoggerConfig) (Logger, error) {
	if config.Underlying == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Underlying must not be empty", config)
	}

	l := &activationKeyLogger{
		underlying: config.Underlying,

		activationKeys: config.ActivationKeys,
	}

	return l, nil
}

func (l *activationKeyLogger) Log(keyVals ...interface{}) error {
	activated, err := shouldActivate(l.activationKeys, keyVals)
	if err != nil {
		return microerror.Mask(err)
	}

	if activated {
		return l.underlying.Log(keyVals...)
	}

	return nil
}

func (l *activationKeyLogger) LogCtx(ctx context.Context, keyVals ...interface{}) error {
	activated, err := shouldActivate(l.activationKeys, keyVals)
	if err != nil {
		return microerror.Mask(err)
	}

	if activated {
		return l.underlying.LogCtx(ctx, keyVals...)
	}

	return nil
}

func (l *activationKeyLogger) With(keyVals ...interface{}) Logger {
	return l.underlying.With(keyVals...)
}

func containsString(keyVals []interface{}, activationKey string) bool {
	for i := 0; i < len(keyVals); i += 2 {
		s, ok := keyVals[i].(string)
		if ok && s == activationKey {
			return true
		}
	}

	return false
}

func isLevelAllowed(keyVals []interface{}, activationKey string) bool {
	activationKeyLevel, ok := levelMapping[activationKey]
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

		return activationKeyLevel >= keyValsLevel
	}

	return false
}

func shouldActivate(activationKeys []string, keyVals []interface{}) (bool, error) {
	var activationCount int

	for _, activationKey := range activationKeys {
		if containsString(keyVals, activationKey) {
			activationCount++
			continue
		}

		if isLevelAllowed(keyVals, activationKey) {
			activationCount++
			continue
		}
	}

	if len(activationKeys) != 0 && len(activationKeys) == activationCount {
		return true, nil
	}

	return false, nil
}
