package micrologger

import (
	"testing"
)

func Test_ActivationKeyLogger_shouldActivate(t *testing.T) {
	testCases := []struct {
		Activations    map[string]string
		KeyVals        []interface{}
		ExpectedResult bool
	}{
		// Case 0, zero value input results into false, because logging should not
		// be activated in case no match exists, even if the input is empty.
		{
			Activations:    nil,
			KeyVals:        nil,
			ExpectedResult: false,
		},

		// Case 1, same as 0 but with empty lists instead of zero values.
		{
			Activations:    map[string]string{},
			KeyVals:        []interface{}{},
			ExpectedResult: false,
		},

		// Case 2, a given activation key not matching any keyVals results into
		// false.
		{
			Activations: map[string]string{
				"foo": "bar",
			},
			KeyVals:        nil,
			ExpectedResult: false,
		},

		// Case 3, same as 2 but with different activation keys.
		// false.
		{
			Activations: map[string]string{
				"foo": "bar",
				"bar": "foo",
				"baz": "foo",
			},
			KeyVals:        nil,
			ExpectedResult: false,
		},

		// Case 4, same as 2 but with given keyVals which still do not match.
		{
			Activations: map[string]string{
				"foo": "bar",
			},
			KeyVals: []interface{}{
				"test",
				3,
				"key",
				"val",
			},
			ExpectedResult: false,
		},

		// Case 5, same as 4 but with different activation keys.
		{
			Activations: map[string]string{
				"foo": "bar",
				"bar": "foo",
				"baz": "foo",
			},
			KeyVals: []interface{}{
				"test",
				3,
				"key",
				"val",
			},
			ExpectedResult: false,
		},

		// Case 6, a given activation key matching any keyVals results into true.
		{
			Activations: map[string]string{
				"test": "*",
			},
			KeyVals: []interface{}{
				"test",
				3,
				"key",
				"val",
			},
			ExpectedResult: true,
		},

		// Case 7, same as 6 but with different matching activation keys.
		{
			Activations: map[string]string{
				"test": "*",
				"key":  "*",
			},
			KeyVals: []interface{}{
				"test",
				3,
				"key",
				"val",
			},
			ExpectedResult: true,
		},

		// Case 8, activation keys matching values of the keyVals result in false,
		// because we only want to activate on matching keys.
		{
			Activations: map[string]string{
				"val": "*",
			},
			KeyVals: []interface{}{
				"key",
				"val",
			},
			ExpectedResult: false,
		},

		// Case 9, activation keys matching keys of the keyVals still result in true
		// even if values match as well.
		{
			Activations: map[string]string{
				"key": "*",
				"val": "*",
			},
			KeyVals: []interface{}{
				"key",
				"val",
				"val",
				"key",
			},
			ExpectedResult: true,
		},

		// Case 10, activation keys must all match in order to result in true.
		{
			Activations: map[string]string{
				"foo": "*",
				"bar": "*",
				"baz": "*",
			},
			KeyVals: []interface{}{
				"foo",
				"val",
				"bar",
				"val",
				"baz",
				"val",
			},
			ExpectedResult: true,
		},

		// Case 11, not all activation keys matching results in false.
		{
			Activations: map[string]string{
				"foo": "*",
				"bar": "*",
				"baz": "*",
			},
			KeyVals: []interface{}{
				"foo",
				"val",
				"bar",
				"val",
				"notmatching",
				"val",
			},
			ExpectedResult: false,
		},

		// Case 12, activation keys representing common log levels result in true
		// when matching.
		{
			Activations: map[string]string{
				"info": "*",
			},
			KeyVals: []interface{}{
				"test",
				3,
				"info",
				"val",
			},
			ExpectedResult: true,
		},

		// Case 13, same as 12 but with a different log level.
		{
			Activations: map[string]string{
				"error": "*",
			},
			KeyVals: []interface{}{
				"test",
				3,
				"error",
				"val",
			},
			ExpectedResult: true,
		},

		// Case 14, activation keys representing common log levels result in true
		// when matching lower log levels. The activation key info matches the log
		// level debug because debug is lower than info.
		{
			Activations: map[string]string{
				"info": "*",
			},
			KeyVals: []interface{}{
				"test",
				3,
				"debug",
				"val",
			},
			ExpectedResult: true,
		},

		// Case 15, activation keys representing common log levels result in false
		// when not matching lower log levels. The activation key info does not
		// match the log level warn because warn is higher than info.
		{
			Activations: map[string]string{
				"info": "*",
			},
			KeyVals: []interface{}{
				"test",
				3,
				"warn",
				"val",
			},
			ExpectedResult: false,
		},

		// Case 16, activation keys representing common log levels result in false
		// when not matching lower log levels. The activation key info does not
		// match the log level error because error is higher than info.
		{
			Activations: map[string]string{
				"info": "*",
			},
			KeyVals: []interface{}{
				"test",
				3,
				"error",
				"val",
			},
			ExpectedResult: false,
		},

		// Case 17, ... .
		{
			Activations: map[string]string{
				"level":     "info",
				"verbosity": "3",
			},
			KeyVals: []interface{}{
				"level",
				"info",
				"verbosity",
				"3",
				"message",
				"test",
			},
			ExpectedResult: true,
		},

		// Case 18, ... .
		{
			Activations: map[string]string{
				"level":     "info",
				"verbosity": "3",
			},
			KeyVals: []interface{}{
				"level",
				"info",
				"verbosity",
				"3",
				"message",
				"test",
			},
			ExpectedResult: true,
		},
	}

	for i, tc := range testCases {
		result, err := shouldActivate(tc.Activations, tc.KeyVals)
		if err != nil {
			t.Fatalf("case %d expected %#v got %#v", i, nil, err)
		}

		if result != tc.ExpectedResult {
			t.Fatalf("case %d expected %#v got %#v", i, tc.ExpectedResult, result)
		}
	}
}
