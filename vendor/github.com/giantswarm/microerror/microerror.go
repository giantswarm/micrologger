package microerror

import (
	"fmt"
	"runtime"

	"errors"
)

// Cause is here only for backward compatibility purposes and should not be used.
//
// NOTE: Use Is instead.
func Cause(err error) error {
	e := err
	for e != nil {
		err = e
		e = errors.Unwrap(err)
	}

	return err
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func Maskf(err Error, f string, v ...interface{}) error {
	annotatedErr := annotatedError{
		annotation: fmt.Sprintf(f, v...),
		underlying: err,
	}

	return mask(annotatedErr)
}

func Mask(err error) error {
	if err == nil {
		return nil
	}

	return mask(err)
}

func mask(err error) error {
	_, file, line, _ := runtime.Caller(2)

	return stackedError{
		stackEntry: StackEntry{
			File: file,
			Line: line,
		},
		underlying: err,
	}
}

// New is here only for backward compatibility purposes and should not be used.
//
// NOTE: Use struct initialization literal for Error struct instead.
func New(kind string) Error {
	return Error{
		Kind: kind,
	}
}
