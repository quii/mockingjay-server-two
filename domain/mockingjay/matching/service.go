package matching

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/quii/mockingjay-server-two/domain/crud"
	"github.com/quii/mockingjay-server-two/domain/mockingjay"
)

type MockingjayStubServerService struct {
	endpoints    crud.CRUDesque[uuid.UUID, mockingjay.Endpoint]
	matchReports crud.CRUDesque[uuid.UUID, Report]
}

func NewMockingjayStubServerService(endpoints mockingjay.Endpoints) (*MockingjayStubServerService, error) {
	endpointCRUD := crud.NewCRUD[uuid.UUID, mockingjay.Endpoint]()
	for _, endpoint := range endpoints {
		if err := endpointCRUD.Create(endpoint.ID, endpoint); err != nil {
			return nil, err
		}
	}
	return &MockingjayStubServerService{endpoints: endpointCRUD, matchReports: crud.NewCRUD[uuid.UUID, Report]()}, nil
}

func (m *MockingjayStubServerService) Reports() crud.CRUDesque[uuid.UUID, Report] {
	return m.matchReports
}

func (m *MockingjayStubServerService) Endpoints() crud.CRUDesque[uuid.UUID, mockingjay.Endpoint] {
	return m.endpoints
}

func (m *MockingjayStubServerService) CreateMatchReport(r *http.Request) (Report, error) {
	endpoints, err := m.endpoints.GetAll()
	if err != nil {
		return Report{}, err
	}
	matchReport := NewReport(r, endpoints)
	if err := m.matchReports.Create(matchReport.ID, matchReport); err != nil {
		return Report{}, err
	}
	return matchReport, nil
}
