package micrologger

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"testing"
	"unicode"

	"github.com/google/go-cmp/cmp"

	"github.com/giantswarm/micrologger/loggermeta"
)

var update = flag.Bool("update", false, "update .golden files")

// Test_MicroLogger tests MicroLogger output.
//
// It uses golden file as reference and when changes to template are
// intentional, they can be updated by providing -update flag for go test.
//
//	go test . -run Test_MicroLogger -update
//
func Test_MicroLogger(t *testing.T) {
	testCases := []struct {
		name              string
		inputCtxKeyValues map[string]string
		inputLogKeyVals   []interface{}
		inputWithKeyVals  []interface{}
	}{
		{
			name:              "case 0",
			inputCtxKeyValues: map[string]string{},
			inputLogKeyVals: []interface{}{
				"foo", "bar",
			},
			inputWithKeyVals: []interface{}{},
		},
		{
			name: "case 1",
			inputCtxKeyValues: map[string]string{
				"baz": "zap",
			},
			inputLogKeyVals: []interface{}{
				"foo", "bar",
			},
			inputWithKeyVals: []interface{}{},
		},
		{
			name:              "case 2",
			inputCtxKeyValues: map[string]string{},
			inputLogKeyVals: []interface{}{
				"foo", "bar",
			},
			inputWithKeyVals: []interface{}{
				"baz", "zap",
			},
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var err error

			w := &bytes.Buffer{}

			var log Logger
			{
				c := Config{
					IOWriter: w,
					TimestampFormatter: func() interface{} {
						return "2019-10-08T20:04:13.490819+00:00"
					},
				}

				log, err = New(c)
				if err != nil {
					t.Fatalf("err = %v, want %v", err, nil)
				}
			}

			if len(tc.inputWithKeyVals) > 0 {
				log = log.With(tc.inputWithKeyVals...)
			}
			if len(tc.inputCtxKeyValues) == 0 {
				log.Log(tc.inputLogKeyVals...)
			} else {
				meta := loggermeta.New()
				meta.KeyVals = tc.inputCtxKeyValues

				ctx := loggermeta.NewContext(context.Background(), meta)

				log.LogCtx(ctx, tc.inputLogKeyVals...)
			}

			var actual []byte
			{
				// Don't flush on purpose. Logs should be
				// flushed right after they are logged.
				wCopy := []byte(w.String())
				w.Reset()
				err := json.Indent(w, wCopy, "", "\t")
				if err != nil {
					t.Fatalf("err = %v, want %v", err, nil)
				}
				actual = w.Bytes()
			}

			golden := filepath.Join("testdata", normalizeToFileName(tc.name)+".golden")
			if *update {
				ioutil.WriteFile(golden, actual, 0644)
			}

			expected, err := ioutil.ReadFile(golden)
			if err != nil {
				t.Fatal(err)
			}

			if !bytes.Equal(actual, expected) {
				t.Fatalf("\n\n%s\n", cmp.Diff(actual, expected))
			}
		})
	}
}

// normalizeToFileName converts all non-digit, non-letter runes in input string
// to dash ('-'). Coalesces multiple dashes into one.
func normalizeToFileName(s string) string {
	var result []rune
	for _, r := range []rune(s) {
		if unicode.IsDigit(r) || unicode.IsLetter(r) {
			result = append(result, r)
		} else {
			l := len(result)
			if l > 0 && result[l-1] != '-' {
				result = append(result, rune('-'))
			}
		}
	}
	return string(result)
}
