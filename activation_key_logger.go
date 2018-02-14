package micrologger

import (
	"context"

	"github.com/giantswarm/microerror"
)

type ActivationKeyConfig struct {
	Underlying Logger

	ActivationKeys []string
}

type activationKeyLogger struct {
	underlying Logger

	activationKeys []string
}

func NewActivationKey(config ActivationKeyConfig) (Logger, error) {
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
	activated, err := shouldActivate(l.activationKeys, keyVals...)
	if err != nil {
		return microerror.Mask(err)
	}

	if activated {
		return l.underlying.Log(keyVals...)
	}

	return nil
}

func (l *activationKeyLogger) LogCtx(ctx context.Context, keyVals ...interface{}) error {
	activated, err := shouldActivate(l.activationKeys, keyVals...)
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

// TODO implement properly.
func shouldActivate(activationKeys []string, keyVals ...interface{}) (bool, error) {
	return false, nil
}
