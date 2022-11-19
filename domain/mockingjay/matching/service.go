package matching

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/quii/mockingjay-server-two/domain/crud"
	"github.com/quii/mockingjay-server-two/domain/mockingjay"
)

type MockingjayStubServerService struct {
	endpoints    crud.CRUD[uuid.UUID, mockingjay.Endpoint]
	matchReports crud.CRUD[uuid.UUID, Report]
}

func NewMockingjayStubServerService(endpoints mockingjay.Endpoints) (*MockingjayStubServerService, error) {
	if err := endpoints.Compile(); err != nil {
		return nil, err
	}
	return &MockingjayStubServerService{endpoints: mockingjay.NewEndpointCRUD(endpoints), matchReports: NewReportCRUD()}, nil
}

func (m *MockingjayStubServerService) Reports() crud.CRUD[uuid.UUID, Report] {
	return m.matchReports
}

func (m *MockingjayStubServerService) Endpoints() crud.CRUD[uuid.UUID, mockingjay.Endpoint] {
	return m.endpoints
}

func (m *MockingjayStubServerService) Reset() {
	m.endpoints = mockingjay.NewEndpointCRUD(nil)
	m.matchReports = NewReportCRUD()
}

func (m *MockingjayStubServerService) CreateMatchReport(r *http.Request) (Report, error) {
	endpoints, err := m.endpoints.GetAll()
	if err != nil {
		return Report{}, err
	}
	matchReport := NewReport(r, endpoints)
	if err := m.matchReports.Create(matchReport); err != nil {
		return Report{}, err
	}
	return matchReport, nil
}
