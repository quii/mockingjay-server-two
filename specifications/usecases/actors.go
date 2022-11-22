package usecases

import (
	"github.com/google/uuid"
	"github.com/quii/mockingjay-server-two/domain/mockingjay/http"
	"github.com/quii/mockingjay-server-two/domain/mockingjay/matching"
)

type Admin interface {
	GetReports() ([]matching.Report, error)
	DeleteReports() error

	AddEndpoints(endpoints ...http.Endpoint) error
	GetEndpoints() (http.Endpoints, error)
	DeleteEndpoint(uuid uuid.UUID) error
	DeleteEndpoints() error
}

type Client interface {
	Send(request http.Request) (http.Response, matching.Report, error)
	//CheckEndpoints() ([]contract.Report, error) - wip
}
