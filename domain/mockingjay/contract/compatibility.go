package contract

import (
	"github.com/quii/mockingjay-server-two/domain/mockingjay/matching"
	"github.com/quii/mockingjay-server-two/domain/mockingjay/stub"
)

func IsResponseCompatible(got, want stub.Response) bool {
	return want.Status == got.Status &&
		matchBodies(got, want) &&
		matching.MatchHeaders(want.Headers, got.Headers)
}

func matchBodies(got stub.Response, want stub.Response) bool {
	return want.Body == got.Body
}
