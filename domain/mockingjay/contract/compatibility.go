package contract

import (
	"encoding/json"

	"github.com/quii/mockingjay-server-two/domain/mockingjay/contract/jsonequaliser"
	"github.com/quii/mockingjay-server-two/domain/mockingjay/matching"
	"github.com/quii/mockingjay-server-two/domain/mockingjay/stub"
)

func IsResponseCompatible(got, want stub.Response) (bool, ErrCompatibilityProblems) {
	errors := ErrCompatibilityProblems{}

	if want.Status != got.Status {
		errors.Errors = append(errors.Errors, ErrorWrongStatus{
			Got:  got.Status,
			Want: want.Status,
		}.Error())
	}

	bodyMatched, bodyErrs := matchBodies(got.Body, want.Body)
	if !bodyMatched {
		errors.Errors = append(errors.Errors, ErrJSONBodyMismatch{Problems: bodyErrs}.Error())
	}

	if !matching.MatchHeaders(want.Headers, got.Headers) {
		errors.Errors = append(errors.Errors, ErrHeadersIncorrect.Error())
	}

	return len(errors.Errors) == 0, errors
}

func matchBodies(got, want string) (bool, map[string]string) {
	if isJSON(want) {
		messages, err := jsonequaliser.IsCompatible(want, got)

		if err != nil {
			return false, map[string]string{"error": err.Error()}
		}
		if len(messages) > 0 {
			return false, messages
		}

		return true, nil
	}

	return want == got, nil
}

func isJSON(s string) bool {
	var js interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}
