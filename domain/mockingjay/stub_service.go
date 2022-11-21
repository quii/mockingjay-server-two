package mockingjay

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/quii/mockingjay-server-two/domain/crud"
	http2 "github.com/quii/mockingjay-server-two/domain/mockingjay/http"
	"github.com/quii/mockingjay-server-two/domain/mockingjay/matching"
)

type StubService struct {
	endpoints    crud.CRUDesque[uuid.UUID, http2.Endpoint]
	matchReports crud.CRUDesque[uuid.UUID, matching.Report]
}

func NewStubService(endpoints http2.Endpoints) (*StubService, error) {
	reportsCRUD := crud.New[uuid.UUID, matching.Report](matching.SortReport)
	endpointCRUD := crud.New[uuid.UUID, http2.Endpoint](http2.SortEndpoint)

	for _, endpoint := range endpoints {
		if err := endpointCRUD.Create(endpoint.ID, endpoint); err != nil {
			return nil, err
		}
	}
	return &StubService{
		endpoints:    endpointCRUD,
		matchReports: reportsCRUD,
	}, nil
}

func (m *StubService) Reports() crud.CRUDesque[uuid.UUID, matching.Report] {
	return m.matchReports
}

func (m *StubService) Endpoints() crud.CRUDesque[uuid.UUID, http2.Endpoint] {
	return m.endpoints
}

func (m *StubService) CreateMatchReport(r *http.Request) (matching.Report, error) {
	endpoints, err := m.endpoints.GetAll()
	if err != nil {
		return matching.Report{}, err
	}
	matchReport := matching.NewReport(r, endpoints)
	if err := m.matchReports.Create(matchReport.ID, matchReport); err != nil {
		return matching.Report{}, err
	}
	return matchReport, nil
}
