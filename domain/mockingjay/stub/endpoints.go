package stub

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
		Ignore    bool   `json:"ignore"`
	}
)

func SortEndpoint(a, b Endpoint) bool {
	return a.LoadedAt.Before(b.LoadedAt)
}

func (e *Endpoint) Compile() error {
	if err := e.Request.compile(); err != nil {
		return err
	}
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	if e.LoadedAt.IsZero() {
		e.LoadedAt = time.Now()
	}
	e.Request.Headers.compile()
	e.Response.Headers.compile()
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
