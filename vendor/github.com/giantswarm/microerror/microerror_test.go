package microerror

import (
	"errors"
	"strconv"
	"testing"
)

func Test_Cause(t *testing.T) {
	var testCauseMicroError = Error{
		Kind: "testCauseErrorB",
	}
	var testCauseErrorsNewError = errors.New("test cause error A")

	testCases := []struct {
		name               string
		inputErrorFunc     func() error
		expectedCauseError error
	}{
		{
			name: "case 0: nil",
			inputErrorFunc: func() error {
				return nil
			},
			expectedCauseError: nil,
		},
		{
			name: "case 1: no masking error=microerror.Error",
			inputErrorFunc: func() error {
				err := testCauseMicroError
				return err
			},
			expectedCauseError: testCauseMicroError,
		},
		{
			name: "case 2: no masking error=errors.New",
			inputErrorFunc: func() error {
				err := testCauseErrorsNewError
				return err
			},
			expectedCauseError: testCauseErrorsNewError,
		},
		{
			name: "case 3: Maskf depth=1 error=microerror.Error",
			inputErrorFunc: func() error {
				err := Maskf(testCauseMicroError, "test annotation")
				return err
			},
			expectedCauseError: testCauseMicroError,
		},
		{
			name: "case 4: Maskf depth=3 error=microerror.Error",
			inputErrorFunc: func() error {
				err := Maskf(testCauseMicroError, "test annotation")
				err = Mask(err)
				err = Mask(err)
				err = Mask(err)
				return err
			},
			expectedCauseError: testCauseMicroError,
		},
		{
			name: "case 5: Mask depth=1 error=microerror.Error",
			inputErrorFunc: func() error {
				err := Mask(testCauseErrorsNewError)
				return err
			},
			expectedCauseError: testCauseErrorsNewError,
		},
		{
			name: "case 6: Mask depth=3 error=microerror.Error",
			inputErrorFunc: func() error {
				err := Mask(testCauseErrorsNewError)
				err = Mask(err)
				err = Mask(err)
				return err
			},
			expectedCauseError: testCauseErrorsNewError,
		},
		{
			name: "case 7: Mask depth=1 error=microerror.Error",
			inputErrorFunc: func() error {
				err := Mask(testCauseErrorsNewError)
				return err
			},
			expectedCauseError: testCauseErrorsNewError,
		},
		{
			name: "case 8: Mask depth=3 error=errors.New",
			inputErrorFunc: func() error {
				err := Mask(testCauseErrorsNewError)
				err = Mask(err)
				err = Mask(err)
				return err
			},
			expectedCauseError: testCauseErrorsNewError,
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			err := tc.inputErrorFunc()
			cause := Cause(err)
			if cause != tc.expectedCauseError {
				t.Errorf("err = %#v, want %#v", cause, tc.expectedCauseError)
			}
		})
	}
}

func Test_Mask_Nil(t *testing.T) {
	err := Mask(nil)
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

var testMicroErr = Error{
	Desc: "test-desc",
	Docs: "test-docs",
	Kind: "testKind",
}

func Test_Mask_Error(t *testing.T) {
	testCases := []struct {
		name           string
		inputErrorFunc func() error
		expectedError  string
	}{
		{
			name: "case 0: Maskf depth=1 error=microerror.Error",
			inputErrorFunc: func() error {
				err := Maskf(testMicroErr, "test annotation")
				return err
			},
			expectedError: "test kind: test annotation",
		},
		{
			name: "case 1: Maskf depth=3 error=microerror.Error",
			inputErrorFunc: func() error {
				err := Maskf(testMicroErr, "test annotation")
				err = Mask(err)
				err = Mask(err)
				return err
			},
			expectedError: "test kind: test annotation",
		},
		{
			name: "case 2: Maskf depth=3 error=microerror.Error & empty annotation",
			inputErrorFunc: func() error {
				err := Maskf(testMicroErr, "")
				err = Mask(err)
				err = Mask(err)
				return err
			},
			expectedError: "test kind",
		},
		{
			name: "case 3: Mask depth=1 error=microerror.Error",
			inputErrorFunc: func() error {
				err := Mask(testMicroErr)
				return err
			},
			expectedError: "test kind",
		},
		{
			name: "case 4: Mask depth=3 error=microerror.Error",
			inputErrorFunc: func() error {
				err := Mask(testMicroErr)
				err = Mask(err)
				err = Mask(err)
				return err
			},
			expectedError: "test kind",
		},
		{
			name: "case 5: Mask depth=1 error=errors.New",
			inputErrorFunc: func() error {
				err := Mask(errors.New("test error"))
				return err
			},
			expectedError: "test error",
		},
		{
			name: "case 6: Mask depth=3 error=errors.New",
			inputErrorFunc: func() error {
				err := Mask(errors.New("test error"))
				err = Mask(err)
				err = Mask(err)
				return err
			},
			expectedError: "test error",
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			err := tc.inputErrorFunc()
			if err.Error() != tc.expectedError {
				t.Errorf("err.Error() = %#q, want %#q", err.Error(), tc.expectedError)
			}
		})
	}
}
