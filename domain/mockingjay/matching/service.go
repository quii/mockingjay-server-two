package matching

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/quii/mockingjay-server-two/domain/mockingjay"
)

type MockingjayStubServerService struct {
	endpoints    mockingjay.Endpoints
	matchReports map[uuid.UUID]Report
}

func NewMockingjayStubServerService(endpoints mockingjay.Endpoints) *MockingjayStubServerService {
	return &MockingjayStubServerService{endpoints: endpoints, matchReports: make(map[uuid.UUID]Report)}
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

func (m *MockingjayStubServerService) PutEndpoints(e mockingjay.Endpoints) error {
	if err := e.Compile(); err != nil {
		return err
	}
	m.endpoints = e
	return nil
}

func (m *MockingjayStubServerService) GetEndpoints() mockingjay.Endpoints {
	return m.endpoints
}
