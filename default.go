package micrologger

import (
	"os"
	"time"

	kitlog "github.com/go-kit/log"
)

var DefaultCaller = newCallerFunc(0)

var DefaultIOWriter = os.Stdout

var DefaultTimestampFormatter kitlog.Valuer = func() interface{} {
	return time.Now().UTC().Format("2006-01-02T15:04:05.999999-07:00")
}
