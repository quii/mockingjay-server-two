package mockingjay

import (
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/exp/slices"
)

type EndpointCRUD struct {
	endpoints Endpoints
}

func NewEndpointCRUD(endpoints Endpoints) *EndpointCRUD {
	return &EndpointCRUD{endpoints: endpoints}
}

func (e *EndpointCRUD) GetAll() ([]Endpoint, error) {
	return e.endpoints, nil
}

func (e *EndpointCRUD) GetByID(id uuid.UUID) (Endpoint, bool, error) {
	i := slices.IndexFunc(e.endpoints, func(endpoint Endpoint) bool {
		return endpoint.ID == id
	})

	if i == -1 {
		return Endpoint{}, false, nil
	}

	return e.endpoints[i], true, nil
}

func (e *EndpointCRUD) Create(endpoint Endpoint) error {
	e.endpoints = append(e.endpoints, endpoint)
	return nil
}

func (e *EndpointCRUD) Delete(id uuid.UUID) error {
	i := slices.IndexFunc(e.endpoints, func(endpoint Endpoint) bool {
		return endpoint.ID == id
	})
	if i == -1 {
		return fmt.Errorf("no such id %s to delete", id.String())
	}
	e.endpoints = slices.Delete(e.endpoints, i, i+1)
	return nil
}
