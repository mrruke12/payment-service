package enum

import "errors"

var (
	ErrInvalidKey error = errors.New("the given key is not present in the set")
)
