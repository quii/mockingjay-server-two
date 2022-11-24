package contract

import "github.com/quii/mockingjay-server-two/domain/mockingjay/stub"

func IsResponseCompatible(got, want stub.Response) bool {
	return got.Status == want.Status && got.Body == want.Body //todo: match headers, re-use the other code probs
}
