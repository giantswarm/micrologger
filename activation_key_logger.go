package micrologger

import (
	"context"

	"github.com/giantswarm/microerror"
)

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

func shouldActivate(activationKeys []string, keyVals []interface{}) (bool, error) {
	for _, k := range activationKeys {
		if containsString(keyVals, k) {
			return true, nil
		}
	}

	return false, nil
}

func containsString(list []interface{}, key string) bool {
	for _, v := range list {
		s, ok := v.(string)
		if ok && s == key {
			return true
		}
	}

	return false
}
