package micrologger

import (
	"fmt"

	kitlog "github.com/go-kit/log"
	"github.com/go-stack/stack"
)

func newCallerFunc(skip int) kitlog.Valuer {
	return func() interface{} {
		return fmt.Sprintf("%+v", stack.Caller(5+skip))
	}
}
