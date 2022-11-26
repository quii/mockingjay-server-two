package usecases

import (
	"testing"
)

func RunAllAgainst[T any](t *testing.T, inputs []T, useCase UseCase[T]) {
	t.Helper()
	for _, input := range inputs {
		useCase.Test(t, input)
	}
}

type UseCase[T any] interface {
	Test(t *testing.T, thing T)
}
