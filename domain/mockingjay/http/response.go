package http

type Response struct {
	Status  int     `json:"status,omitempty"`
	Body    string  `json:"body,omitempty"`
	Headers Headers `json:"headers,omitempty"`
}
