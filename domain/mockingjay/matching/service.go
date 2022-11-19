package matching

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/quii/mockingjay-server-two/domain/mockingjay"
	"golang.org/x/exp/slices"
)

type MockingjayStubServerService struct {
	endpoints    mockingjay.Endpoints
	matchReports map[uuid.UUID]Report
}

func NewMockingjayStubServerService(endpoints mockingjay.Endpoints) (*MockingjayStubServerService, error) {
	if err := endpoints.Compile(); err != nil {
		return nil, err
	}
	return &MockingjayStubServerService{endpoints: endpoints, matchReports: make(map[uuid.UUID]Report)}, nil
}

func (m *MockingjayStubServerService) Reset() {
	m.endpoints = mockingjay.Endpoints{}
	m.matchReports = make(map[uuid.UUID]Report)
}

func (m *MockingjayStubServerService) GetMatchReport(r *http.Request) Report {
	matchReport := NewReport(r, m.endpoints)
	m.matchReports[matchReport.ID] = matchReport
	return matchReport
}

func (m *MockingjayStubServerService) GetReports() Reports {
	var reports Reports
	for _, report := range m.matchReports {
		reports = append(reports, report)
	}
	reports.Sort()
	return reports
}

func (m *MockingjayStubServerService) GetReport(id uuid.UUID) (Report, bool) {
	report, exists := m.matchReports[id]
	return report, exists
}

func (m *MockingjayStubServerService) AddEndpoint(e mockingjay.Endpoint) error {
	m.endpoints = append(m.endpoints, e)
	return nil
}

func (m *MockingjayStubServerService) GetEndpoints() mockingjay.Endpoints {
	return m.endpoints
}

func (m *MockingjayStubServerService) DeleteEndpoint(id uuid.UUID) error {
	i := slices.IndexFunc(m.endpoints, func(endpoint mockingjay.Endpoint) bool {
		return endpoint.ID == id
	})
	if i == -1 {
		return fmt.Errorf("no such id %s to delete", id.String())
	}
	m.endpoints = slices.Delete(m.endpoints, i, i+1)
	return nil
}
