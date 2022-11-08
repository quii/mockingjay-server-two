package contract

import "github.com/quii/mockingjay-server-two/domain/mockingjay"

type Report struct {
	mockingjay.Endpoint
	ResponseFromDownstream mockingjay.Response
	Passed                 bool `json:"passed"`
}
