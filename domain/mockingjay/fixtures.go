package mockingjay

type (
	RequestDescription struct {
		Description string  `json:"description,omitempty"`
		Request     Request `json:"request"`
	}

	TestFixture struct {
		Endpoint            Endpoint             `json:"endpoint"`
		MatchingRequests    []RequestDescription `json:"matchingRequests,omitempty"`
		NonMatchingRequests []RequestDescription `json:"nonMatchingRequests,omitempty"`
	}
)
