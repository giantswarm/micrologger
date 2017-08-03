package micrologger

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestWith(t *testing.T) {
	var err error

	out := new(bytes.Buffer)

	var log Logger
	{
		c := DefaultConfig()
		{
			c.IOWriter = out
		}
		log, err = New(c)
		if err != nil {
			panic(err)
		}
	}

	var (
		field       = "ctxField"
		wfieldValue = "test ctx field value"

		parentLog = log
		childLog  = log.With(field, wfieldValue)
	)

	// Make sure caller (old field) and added contextual field are logged.
	{
		wfieldValue := "test ctx field value"

		out.Reset()
		childLog.Log("msg", "whats up bro?")

		var got map[string]string
		json.Unmarshal(out.Bytes(), &got)

		wcaller := "github.com/giantswarm/micrologger/logger_test.go:39"
		caller, ok := got["caller"]
		if !ok {
			t.Errorf("expected caller key")
		}
		if caller != wcaller {
			t.Errorf("want caller %s, got %s", wcaller, caller)
		}

		fieldValue, ok := got[field]
		if !ok {
			t.Errorf("want set field %s", field)
		}
		if fieldValue != wfieldValue {
			t.Errorf("want fieldValue %s, got %s", wfieldValue, fieldValue)
		}
	}

	// Make sure parent logger remained unchanged.
	{
		out.Reset()
		parentLog.Log("msg", "how are you?")

		var got map[string]string
		json.Unmarshal(out.Bytes(), &got)

		wcaller := "github.com/giantswarm/micrologger/logger_test.go:65"
		caller, ok := got["caller"]
		if !ok {
			t.Errorf("expected caller key")
		}
		if caller != wcaller {
			t.Errorf("want caller %s, got %s", wcaller, caller)
		}

		fieldValue, ok := got[field]
		if ok {
			t.Errorf("want unset field %s, got %s", field, fieldValue)
		}
	}
}
