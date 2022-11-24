package mockingjay

import "github.com/quii/mockingjay-server-two/domain/mockingjay/stub"

type (
	RequestDescription struct {
		Description string       `json:"description,omitempty"`
		Request     stub.Request `json:"request"`
	}

	TestFixture struct {
		Endpoint            stub.Endpoint        `json:"endpoint"`
		MatchingRequests    []RequestDescription `json:"matchingRequests,omitempty"`
		NonMatchingRequests []RequestDescription `json:"nonMatchingRequests,omitempty"`
	}
)
