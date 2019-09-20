package microerror

import (
	"encoding/json"
	"errors"
)

// JSON prints the error with enriched information in JSON format. Enriched
// information includes:
//
//	- All fields from Error type.
//	- Error stack.
//
// The rendered JSON can be unmarshalled with JSONError type.
func JSON(err error) string {
	if !errors.As(err, &Error{}) && !errors.As(err, &stackedError{}) {
		err = annotatedError{
			annotation: err.Error(),
			underlying: Error{
				Kind: kindUnknown,
			},
		}
	}

	bytes, err := json.Marshal(err)
	if err != nil {
		panic(err.Error())
	}

	return string(bytes)
}
