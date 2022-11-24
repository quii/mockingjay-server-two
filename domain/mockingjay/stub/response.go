package stub

import (
	"io"
	"net/http"
)

type Response struct {
	Status  int     `json:"status,omitempty"`
	Body    string  `json:"body,omitempty"`
	Headers Headers `json:"headers,omitempty"`
}

func NewResponseFromHTTP(res *http.Response) Response {
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	return Response{
		Status:  res.StatusCode,
		Body:    string(body),
		Headers: Headers(res.Header),
	}
}
