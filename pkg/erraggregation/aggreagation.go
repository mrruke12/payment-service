package erraggregation

import (
	"errors"
	"fmt"
	"reflect"
)

// Join joins errors and strings in the given order.
//
// It is a syntax sugar over [fmt.Errorf]: errors are wrapped with %w,
// strings are formatted with %s. Values of other types are skipped.
//
// Valid arg types: string, error.
//
// Example:
//
//	err := Join("failed to load", ErrNotFound, "for user", userID)
//	// => "failed to load: not found: for user 42"
func Join(args ...any) error {
	var acc error

	for _, arg := range args {
		switch v := arg.(type) {
		case error:
			acc = joinError(acc, v)
		case string:
			acc = joinString(acc, v)
		default:
			if reflect.TypeOf(arg).Kind() == reflect.String { // handle ~string types
				acc = joinString(acc, reflect.ValueOf(arg).String())
			}
		}
	}

	if acc == nil {
		return errors.New("no args of type string | error passed to apperr.Join")
	}

	return acc
}

// joinError is a sugar over fmt.Errorf("%w: %w", err, sub).
func joinError(err error, sub error) error {
	if err == nil {
		return sub
	}

	return fmt.Errorf("%w: %w", err, sub)
}

// joinString is a sugar over fmt.Errorf("%w: %s", err, sub).
func joinString(err error, sub string) error {
	if err == nil {
		return errors.New(sub)
	}

	return fmt.Errorf("%w: %s", err, sub)
}
