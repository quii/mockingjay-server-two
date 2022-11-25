package collections

func Reduce[A, B any](collection []A, accumulator func(A, B) (B, error), initialValue B) (B, error) {
	var result = initialValue
	var emptyB B
	for _, x := range collection {
		acc, err := accumulator(x, result)
		if err != nil {
			return emptyB, err
		}
		result = acc
	}
	return result, nil
}

func Map[A, B any](collection []A, f func(A) B) []B {
	var result []B
	for _, a := range collection {
		result = append(result, f(a))
	}
	return result
}

func ForAll[A any](collection []A, f func(A) error) error {
	for _, x := range collection {
		if err := f(x); err != nil {
			return err
		}
	}
	return nil
}
