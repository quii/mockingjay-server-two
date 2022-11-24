package usecases

import (
	"github.com/google/uuid"
	"github.com/quii/mockingjay-server-two/domain/mockingjay/contract"
	"github.com/quii/mockingjay-server-two/domain/mockingjay/matching"
	"github.com/quii/mockingjay-server-two/domain/mockingjay/stub"
)

type Admin interface {
	GetReports() ([]matching.Report, error)
	DeleteReports() error

	AddEndpoints(endpoints ...stub.Endpoint) error
	GetEndpoints() (stub.Endpoints, error)
	DeleteEndpoint(uuid uuid.UUID) error
	DeleteEndpoints() error
}

type StubServerClient interface {
	Send(request stub.Request) (stub.Response, matching.Report, error)
}

type ConsumerDrivenContractChecker interface {
	CheckEndpoints() ([]contract.Report, error)
}
