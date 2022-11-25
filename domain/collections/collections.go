package collections

func Reduce[A, B any](collection []A, accumulator func(B, A) B, initialValue B) B {
	var result = initialValue
	for _, x := range collection {
		result = accumulator(result, x)
	}
	return result
}

func Map[A, B any](collection []A, f func(A) B) []B {
	var result []B
	for _, a := range collection {
		result = append(result, f(a))
	}
	return result
}
