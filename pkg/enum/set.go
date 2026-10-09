package enum

import "slices"

// An implementation of an immutable set of string-like values
type Set[T ~string] struct {
	values []T
}

// Initializes a new set
func NewSet[T ~string](values ...T) *Set[T] {
	return &Set[T]{
		values: values,
	}
}

// Reports if the given value is present in the set
func (s *Set[T]) Has(value T) bool {
	return slices.Contains(s.values, value)
}

// Returns a copy of set values
func (s *Set[T]) Values() []T {
	values := make([]T, len(s.values))
	copy(values, s.values)

	return values
}

// Returns a copy of set values as a slice of strings
func (s *Set[T]) Strings() []string {
	values := make([]string, len(s.values))

	for i := range s.values {
		values = append(values, string(s.values[i]))
	}

	return values
}
