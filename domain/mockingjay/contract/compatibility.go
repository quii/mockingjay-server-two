package contract

import (
	"encoding/json"

	"github.com/quii/mockingjay-server-two/domain/mockingjay/contract/jsonequaliser"
	"github.com/quii/mockingjay-server-two/domain/mockingjay/matching"
	"github.com/quii/mockingjay-server-two/domain/mockingjay/stub"
)

func IsResponseCompatible(got, want stub.Response) bool {
	return want.Status == got.Status &&
		matchBodies(got.Body, want.Body) &&
		matching.MatchHeaders(want.Headers, got.Headers)
}

func matchBodies(got, want string) bool {
	if isJSON(want) {
		messages, err := jsonequaliser.IsCompatible(want, got)

		//todo: need to design this to return more details
		if err != nil {
			return false
		}
		if len(messages) > 0 {
			return false
		}

		return true
	}

	return want == got
}

func isJSON(s string) bool {
	var js interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}
