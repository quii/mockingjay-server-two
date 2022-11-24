package contract

import (
	"github.com/quii/mockingjay-server-two/domain/mockingjay/stub"
	"golang.org/x/exp/slices"
)

func IsResponseCompatible(got, want stub.Response) bool {
	return got.Status == want.Status &&
		got.Body == want.Body &&
		checkHeaders(got.Headers, want.Headers)
}

func checkHeaders(got stub.Headers, want stub.Headers) bool {
	for key, values := range want {
		valuesGot, exist := got[key]

		if !exist {
			return false
		}

		for _, valueNeeded := range values {
			if !slices.Contains(valuesGot, valueNeeded) {
				return false
			}
		}
	}

	return true
}
