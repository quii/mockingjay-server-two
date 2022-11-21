package mockingjay

import "github.com/quii/mockingjay-server-two/domain/mockingjay/http"

type (
	RequestDescription struct {
		Description string       `json:"description,omitempty"`
		Request     http.Request `json:"request"`
	}

	TestFixture struct {
		Endpoint            http.Endpoint        `json:"endpoint"`
		MatchingRequests    []RequestDescription `json:"matchingRequests,omitempty"`
		NonMatchingRequests []RequestDescription `json:"nonMatchingRequests,omitempty"`
	}
)
