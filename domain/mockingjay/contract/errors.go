package contract

import (
	"errors"
	"fmt"
	"strings"
)

type (
	ErrCompatibilityProblems struct {
		Errors []string `json:"errors,omitempty"`
	}
	ErrorWrongStatus struct {
		Got  int `json:"got"`
		Want int `json:"want"`
	}
	ErrJSONBodyMismatch struct {
		Problems map[string]string `json:"problems,omitempty"`
	}
)

func (e ErrCompatibilityProblems) Error() string {
	if len(e.Errors) == 0 {
		return ""
	}
	return fmt.Sprintf("failed compatibility check: %s", strings.Join(e.Errors, ", "))
}

var (
	ErrHeadersIncorrect = errors.New("headers are incorrect")
)

func (e ErrorWrongStatus) Error() string {
	return fmt.Sprintf("got status %d, wanted %d", e.Got, e.Want)
}

func (e ErrJSONBodyMismatch) Error() string {
	var fieldsErrors []string
	for k, v := range e.Problems {
		fieldsErrors = append(fieldsErrors, fmt.Sprintf("[%s]: %s", k, v))
	}
	return fmt.Sprintf("json body detected, but was incompatible; %s", strings.Join(fieldsErrors, ", "))
}
