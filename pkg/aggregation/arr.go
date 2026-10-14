package aggregation

func Map[I, O any](values []I, predicate func(value I, index int) O) []O {
	res := make([]O, len(values))

	for i := range values {
		res[i] = predicate(values[i], i)
	}

	return res
}
