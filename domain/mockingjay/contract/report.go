package contract

import (
	"github.com/quii/mockingjay-server-two/domain/mockingjay/stub"
)

type Report struct {
	Endpoint               stub.Endpoint `json:"endpoint"`
	ResponseFromDownstream stub.Response `json:"responseFromDownstream"`
	URL                    string        `json:"URL"`
	Ignore                 bool          `json:"ignore"`
	Errors                 []string      `json:"errors"`
}

func (r Report) Passed() bool {
	return len(r.Errors) == 0
}

func (r Report) PassedOrIgnored() bool {
	return r.Passed() || r.Ignore
}
