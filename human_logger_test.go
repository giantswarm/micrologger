package micrologger

import (
	"bytes"
	"context"
	"testing"

	"github.com/fatih/color"
	"github.com/giantswarm/micrologger/loggermeta"
)

func Test_Logger_LogWithCtxHuman(t *testing.T) {
	var err error

	out := new(bytes.Buffer)
	time := "2006-01-02T15:04:05.999999-07:00"

	var log Logger
	{
		c := Config{
			IOWriter: out,
			Human:    true,
			TimestampFormatter: func() interface{} {
				return time
			},
		}
		log, err = New(c)
		if err != nil {
			t.Fatalf("setting up logger: %#v", err)
		}
	}

	{
		log.LogCtx(context.TODO(), "message", "foo")
		expected := color.GreenString(time) + " INFO foo\n"
		got := string(out.Bytes())

		if expected != got {
			t.Fatalf("expected %v got %v", expected, got)
		}
	}

	var ctx context.Context
	{
		meta := loggermeta.New()
		meta.KeyVals["baz"] = "zap"

		ctx = loggermeta.NewContext(context.Background(), meta)
	}

	{
		out.Reset()
		log.LogCtx(ctx, "message", "bar", "level", "warn")

		expected := color.GreenString(time) + " " + color.YellowString("WARN") + " bar\n"
		got := string(out.Bytes())

		if expected != got {
			t.Fatalf("expected %s got %s", expected, got)
		}
	}
}
