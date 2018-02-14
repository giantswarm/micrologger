package micrologger

import (
	"testing"
)

func Test_ActivationKeyLogger_shouldActivate(t *testing.T) {
	testCases := []struct {
		ActivationKeys []string
		KeyVals        []interface{}
		ExpectedResult bool
	}{
		// Case 0, zero value input results into false, because logging should not
		// be activated in case no match exists, even if the input is empty.
		{
			ActivationKeys: nil,
			KeyVals:        nil,
			ExpectedResult: false,
		},

		// Case 1, same as 0 but with empty lists instead of zero values.
		{
			ActivationKeys: []string{},
			KeyVals:        []interface{}{},
			ExpectedResult: false,
		},

		// Case 2, a given activation key not matching any keyVals results into
		// false.
		{
			ActivationKeys: []string{
				"foo",
			},
			KeyVals:        nil,
			ExpectedResult: false,
		},

		// Case 3, same as 2 but with different activation keys.
		// false.
		{
			ActivationKeys: []string{
				"foo",
				"foo",
				"bar",
				"baz",
			},
			KeyVals:        nil,
			ExpectedResult: false,
		},

		// Case 4, same as 2 but with given keyVals which still do not match.
		{
			ActivationKeys: []string{
				"foo",
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
			ActivationKeys: []string{
				"foo",
				"foo",
				"bar",
				"baz",
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
			ActivationKeys: []string{
				"test",
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
			ActivationKeys: []string{
				"test",
				"key",
			},
			KeyVals: []interface{}{
				"test",
				3,
				"key",
				"val",
			},
			ExpectedResult: true,
		},
	}

	for i, tc := range testCases {
		result, err := shouldActivate(tc.ActivationKeys, tc.KeyVals)
		if err != nil {
			t.Fatalf("case %d expected %#v got %#v", i, nil, err)
		}

		if result != tc.ExpectedResult {
			t.Fatalf("case %d expected %#v got %#v", i, tc.ExpectedResult, result)
		}
	}
}
