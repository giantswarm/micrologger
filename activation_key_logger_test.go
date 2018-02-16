package micrologger

import (
	"testing"
)

func Test_ActivationKeyLogger_shouldActivate_zeroValue(t *testing.T) {
	testCases := []struct {
		Activations    map[string]interface{}
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
			Activations:    map[string]interface{}{},
			KeyVals:        []interface{}{},
			ExpectedResult: false,
		},

		// Case 2, same as 0 but with different input.
		{
			Activations:    nil,
			KeyVals:        []interface{}{},
			ExpectedResult: false,
		},

		// Case 3, same as 0 but with different input.
		{
			Activations:    map[string]interface{}{},
			KeyVals:        nil,
			ExpectedResult: false,
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

func Test_ActivationKeyLogger_shouldActivate_arbitrary(t *testing.T) {
	testCases := []struct {
		Activations    map[string]interface{}
		KeyVals        []interface{}
		ExpectedResult bool
	}{
		// Case 0, a given activation key not matching any keyVals results into
		// false.
		{
			Activations: map[string]interface{}{
				"foo": "bar",
			},
			KeyVals:        nil,
			ExpectedResult: false,
		},

		// Case 1, a given activation key not matching any keyVals results into
		// false.
		{
			Activations: map[string]interface{}{
				"foo": 3,
			},
			KeyVals:        nil,
			ExpectedResult: false,
		},

		// Case 2, same as 2 but with different activation keys.
		// false.
		{
			Activations: map[string]interface{}{
				"foo": "bar",
				"bar": "foo",
				"baz": "foo",
			},
			KeyVals:        nil,
			ExpectedResult: false,
		},

		// Case 3, same as 2 but with given keyVals which still do not match.
		{
			Activations: map[string]interface{}{
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

		// Case 4, same as 4 but with different activation keys.
		{
			Activations: map[string]interface{}{
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

		// Case 5, a given activation key matching any keyVals results into true.
		{
			Activations: map[string]interface{}{
				"test": 3,
			},
			KeyVals: []interface{}{
				"test",
				3,
				"key",
				"val",
			},
			ExpectedResult: true,
		},

		// Case 6, same as 6 but with different matching activation keys.
		{
			Activations: map[string]interface{}{
				"test": 3,
				"key":  "val",
			},
			KeyVals: []interface{}{
				"test",
				3,
				"key",
				"val",
			},
			ExpectedResult: true,
		},

		// Case 7, activation keys matching values of the keyVals result in false,
		// because we only want to activate on matching keys.
		{
			Activations: map[string]interface{}{
				"val": "key",
			},
			KeyVals: []interface{}{
				"key",
				"val",
			},
			ExpectedResult: false,
		},

		// Case 8, activation keys matching keys of the keyVals still result in true
		// even if values match as well.
		{
			Activations: map[string]interface{}{
				"key": "val",
				"val": "key",
			},
			KeyVals: []interface{}{
				"key",
				"val",
				"val",
				"key",
			},
			ExpectedResult: true,
		},

		// Case 9, activation keys must all match in order to result in true.
		{
			Activations: map[string]interface{}{
				"foo": "val",
				"bar": "val",
				"baz": "val",
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

		// Case 10, not all activation keys matching results in false.
		{
			Activations: map[string]interface{}{
				"foo": "val",
				"bar": "val",
				"baz": "val",
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

func Test_ActivationKeyLogger_shouldActivate_level(t *testing.T) {
	testCases := []struct {
		Activations    map[string]interface{}
		KeyVals        []interface{}
		ExpectedResult bool
	}{
		// Case 0, activation keys representing common log levels result in true
		// when matching.
		{
			Activations: map[string]interface{}{
				"level": "info",
			},
			KeyVals: []interface{}{
				"test",
				3,
				"level",
				"info",
			},
			ExpectedResult: true,
		},

		// Case 1, same as 12 but with a different log level.
		{
			Activations: map[string]interface{}{
				"level": "error",
			},
			KeyVals: []interface{}{
				"test",
				3,
				"level",
				"error",
			},
			ExpectedResult: true,
		},

		// Case 2, activation keys representing common log levels result in true
		// when matching lower log levels. The activation key info matches the log
		// level debug because debug is lower than info.
		{
			Activations: map[string]interface{}{
				"level": "info",
			},
			KeyVals: []interface{}{
				"test",
				3,
				"level",
				"debug",
			},
			ExpectedResult: true,
		},

		// Case 3, activation keys representing common log levels result in false
		// when not matching lower log levels. The activation key info does not
		// match the log level warning because warning is higher than info.
		{
			Activations: map[string]interface{}{
				"level": "info",
			},
			KeyVals: []interface{}{
				"test",
				3,
				"level",
				"warning",
			},
			ExpectedResult: false,
		},

		// Case 4, activation keys representing common log levels result in false
		// when not matching lower log levels. The activation key info does not
		// match the log level error because error is higher than info.
		{
			Activations: map[string]interface{}{
				"level": "info",
			},
			KeyVals: []interface{}{
				"test",
				3,
				"level",
				"error",
			},
			ExpectedResult: false,
		},

		// Case 5, ... .
		{
			Activations: map[string]interface{}{
				"level":     "info",
				"verbosity": 3,
			},
			KeyVals: []interface{}{
				"level",
				"info",
				"verbosity",
				3,
				"message",
				"test",
			},
			ExpectedResult: true,
		},

		// Case 6, ... .
		{
			Activations: map[string]interface{}{
				"level":     "info",
				"verbosity": 3,
			},
			KeyVals: []interface{}{
				"level",
				"info",
				"verbosity",
				3,
				"message",
				"test",
			},
			ExpectedResult: true,
		},

		// Case 7, ... independent of verbosity .
		{
			Activations: map[string]interface{}{
				"level": "info",
			},
			KeyVals: []interface{}{
				"level",
				"info",
				"verbosity",
				15,
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

func Test_ActivationKeyLogger_shouldActivate_verbosity(t *testing.T) {
	testCases := []struct {
		Activations    map[string]interface{}
		KeyVals        []interface{}
		ExpectedResult bool
	}{
		// Case 0, ... .
		{
			Activations: map[string]interface{}{
				"verbosity": 3,
			},
			KeyVals: []interface{}{
				"level",
				"info",
				"verbosity",
				3,
				"message",
				"test",
			},
			ExpectedResult: true,
		},

		// Case 1, ... .
		{
			Activations: map[string]interface{}{
				"verbosity": 6,
			},
			KeyVals: []interface{}{
				"level",
				"info",
				"verbosity",
				6,
				"message",
				"test",
			},
			ExpectedResult: true,
		},

		// Case 2, ... .
		{
			Activations: map[string]interface{}{
				"verbosity": 6,
			},
			KeyVals: []interface{}{
				"level",
				"info",
				"verbosity",
				2,
				"message",
				"test",
			},
			ExpectedResult: true,
		},

		// Case 3, ... .
		{
			Activations: map[string]interface{}{
				"verbosity": 6,
			},
			KeyVals: []interface{}{
				"level",
				"info",
				"verbosity",
				12,
				"message",
				"test",
			},
			ExpectedResult: false,
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
