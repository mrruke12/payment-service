package enum

import (
	"fmt"
	"slices"
)

// An implementation of state machine
type StateMachine[T ~string] struct {
	set         *Set[T]
	transitions map[T][]T
}

// Initializes a new state machine with panic on error
func NewStateMachine[T ~string](set *Set[T], transitions map[T][]T) *StateMachine[T] {
	m, err := NewStateMachineSafe(set, transitions)

	if err != nil {
		panic(err)
	}

	return m
}

// Initialize a new state machine without panic on error
func NewStateMachineSafe[T ~string](set *Set[T], transitions map[T][]T) (*StateMachine[T], error) {
	for from, toList := range transitions {
		if !set.Has(from) {
			return nil, fmt.Errorf("could not initialize a state machine for %s: %w", from, ErrInvalidKey)
		}

		for _, to := range toList {
			if !set.Has(to) {
				return nil, fmt.Errorf("could not initialize a state machine for %s: %w", to, ErrInvalidKey)
			}
		}
	}

	return &StateMachine[T]{
		set,
		transitions,
	}, nil
}

// Reports whether the given transition between states is valid
func (m *StateMachine[T]) CanTransition(from, to T) bool {
	if !m.set.Has(from) || !m.set.Has(to) {
		return false
	}

	transitions := m.transitions[from]

	return slices.Contains(transitions, to)
}

// Return the list of transitions for the given state
func (m *StateMachine[T]) NextTransitions(from T) []T {
	if !m.set.Has(from) {
		return []T{}
	}

	transitions := m.transitions[from]

	return slices.Clone(transitions)
}
