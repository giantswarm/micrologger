package microerror

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strconv"
	"testing"
	"unicode"

	"github.com/google/go-cmp/cmp"
)

var update = flag.Bool("update", false, "update resource.golden file")

// Test_JSON tests marshaling errors to JSON.
//
// It uses golden file as reference and when changes to template are
// intentional, they can be updated by providing -update flag for go test.
//
//	go test ./ -run Test_JSON -update
//
func Test_JSON(t *testing.T) {
	testCases := []struct {
		name           string
		inputErrorFunc func() error
	}{
		{
			name: "case 0: error=microerror.Error no masking",
			inputErrorFunc: func() error {
				err := testMicroErr
				return err
			},
		},
		{
			name: "case 1: error=errors.New no masking",
			inputErrorFunc: func() error {
				err := errors.New("test error")
				return err
			},
		},
		{
			name: "case 2: error=microerror.Error depth=1 Maskf",
			inputErrorFunc: func() error {
				err := Maskf(testMicroErr, "test annotation")
				return err
			},
		},
		{
			name: "case 3: error=microerror.Error depth=3 Maskf",
			inputErrorFunc: func() error {
				err := Maskf(testMicroErr, "test annotation")
				err = Mask(err)
				err = Mask(err)
				return err
			},
		},
		{
			name: "case 4: error=microerror.Error depth=1 Mask",
			inputErrorFunc: func() error {
				err := Mask(testMicroErr)
				return err
			},
		},
		{
			name: "case 5: error=errors.New depth=1 Mask",
			inputErrorFunc: func() error {
				err := Mask(errors.New("test error"))
				return err
			},
		},
		{
			name: "case 6: error=microerror.Error depth=3 Mask",
			inputErrorFunc: func() error {
				err := Mask(testMicroErr)
				err = Mask(err)
				// Extra line.
				err = Mask(err)
				return err
			},
		},
		{
			name: "case 7: error=errors.New depth=3 Mask",
			inputErrorFunc: func() error {
				err := Mask(errors.New("test error"))
				err = Mask(err)
				err = Mask(err)
				return err
			},
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			actual := JSON(tc.inputErrorFunc())
			{
				b := &bytes.Buffer{}
				json.Indent(b, []byte(actual), "", "\t")
				actual = b.String()
			}
			// Change paths to avoid prefixes like
			// "/Users/username/go/src/" so this can test can be
			// executed on different machines.
			//
			// E.g: This:
			//
			//	"file":"/Users/username/go/src/github.com/giantswarm/microerror/json_test.go"
			//
			// Should be replaced with:
			//
			//	"file":"--REPLACED--/github.com/giantswarm/microerror/json_test.go"
			//
			{
				r := regexp.MustCompile(`("file"\s*:\s*")[/\w]+(/github.com/giantswarm)`)
				actual = r.ReplaceAllString(actual, "$1--REPLACED--$2")
			}

			var expected string
			{
				golden := filepath.Join("testdata", normalizeToFileName(tc.name)+".golden")
				if *update {
					err := ioutil.WriteFile(golden, []byte(actual), 0644)
					if err != nil {
						t.Fatal(err)
					}
				}

				bytes, err := ioutil.ReadFile(golden)
				if err != nil {
					t.Fatal(err)
				}

				expected = string(bytes)
			}

			if actual != expected {
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
