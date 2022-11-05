package endpoints

type Response struct {
	Status int    `json:"status,omitempty"`
	Body   string `json:"body,omitempty"`
}

type Request struct {
	Method string `json:"method,omitempty"`
	Path   string `json:"path,omitempty"`
}

type Endpoint struct {
	Description string   `json:"description,omitempty"`
	Request     Request  `json:"request"`
	Response    Response `json:"response"`
}
