package contract

import (
	"github.com/quii/mockingjay-server-two/domain/mockingjay/stub"
)

type Report struct {
	Endpoint               stub.Endpoint            `json:"endpoint"`
	ResponseFromDownstream stub.Response            `json:"responseFromDownstream"`
	URL                    string                   `json:"URL"`
	Passed                 bool                     `json:"passed"`
	Errors                 ErrCompatibilityProblems `json:"errors"`
	Ignore                 bool                     `json:"ignore"`
}
