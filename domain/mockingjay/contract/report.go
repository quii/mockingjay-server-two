package contract

import (
	"github.com/quii/mockingjay-server-two/domain/mockingjay/http"
)

type Report struct {
	http.Endpoint
	ResponseFromDownstream http.Response
	Passed                 bool `json:"passed"`
}
