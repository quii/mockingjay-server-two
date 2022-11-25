package contract

import (
	"encoding/json"
	"errors"

	"github.com/quii/mockingjay-server-two/domain/mockingjay/contract/jsonequaliser"
	"github.com/quii/mockingjay-server-two/domain/mockingjay/matching"
	"github.com/quii/mockingjay-server-two/domain/mockingjay/stub"
)

func IsResponseCompatible(got, want stub.Response) []error {
	var errors []error

	if want.Status != got.Status {
		errors = append(errors, ErrorWrongStatus{
			Got:  got.Status,
			Want: want.Status,
		})
	}

	bodyMatched, err := matchBodies(got.Body, want.Body)
	if !bodyMatched {
		errors = append(errors, err)
	}

	if !matching.MatchHeaders(want.Headers, got.Headers) {
		errors = append(errors, ErrHeadersIncorrect)
	}

	return errors
}

func matchBodies(got, want string) (bool, error) {
	if isJSON(want) {
		messages, err := jsonequaliser.IsCompatible(want, got)

		if err != nil {
			return false, err
		}
		if len(messages) > 0 {
			return false, ErrJSONBodyMismatch{Problems: messages}
		}

		return true, nil
	}

	if want != got {
		return false, errors.New("mismatched response bodies")
	}

	return true, nil
}

func isJSON(s string) bool {
	var js interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}
