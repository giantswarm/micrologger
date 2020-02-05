package micrologger

import (
	"github.com/giantswarm/microerror"
)

var invalidConfigError =&
	Kind: "invalidConfigError",
}

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return microerror.Cause(err) == invalidConfigError
}
