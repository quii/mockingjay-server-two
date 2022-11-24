package usecases

import (
	"github.com/google/uuid"
	"github.com/quii/mockingjay-server-two/domain/mockingjay/contract"
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

type StubServerClient interface {
	Send(request http.Request) (http.Response, matching.Report, error)
}

type ConsumerDrivenContractChecker interface {
	CheckEndpoints() ([]contract.Report, error)
}
