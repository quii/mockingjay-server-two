package mockingjay

import (
	"time"

	"github.com/google/uuid"
)

type (
	Endpoints []Endpoint

	Endpoint struct {
		ID          uuid.UUID `json:"ID"`
		Description string    `json:"description,omitempty"`
		Request     Request   `json:"request"`
		Response    Response  `json:"response"`
		CDCs        []CDC
		LoadedAt    time.Time `json:"loadedAt"`
	}

	CDC struct {
		BaseURL   string `json:"baseURL"`
		Retries   int    `json:"retries"`
		TimeoutMS int    `json:"timeoutMS"`
	}

	Response struct {
		Status  int     `json:"status,omitempty"`
		Body    string  `json:"body,omitempty"`
		Headers Headers `json:"headers,omitempty"`
	}

	Headers map[string][]string
)

func (e *Endpoint) Compile() error {
	if err := e.Request.compile(); err != nil {
		return err
	}
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return nil
}

func (e Endpoints) Compile() error {
	for i := range e {
		if err := e[i].Compile(); err != nil {
			return err
		}
	}
	return nil
}
