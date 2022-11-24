package contract

import (
	"github.com/quii/mockingjay-server-two/domain/mockingjay/stub"
)

type Report struct {
	stub.Endpoint
	ResponseFromDownstream stub.Response
	Passed                 bool `json:"passed"`
}
