package mockingjay

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/quii/mockingjay-server-two/domain/crud"
	"github.com/quii/mockingjay-server-two/domain/mockingjay/contract"
	"github.com/quii/mockingjay-server-two/domain/mockingjay/matching"
	"github.com/quii/mockingjay-server-two/domain/mockingjay/stub"
)

type CDCService interface {
	GetReports(endpoint stub.Endpoint) ([]contract.Report, error)
}

type Service struct {
	endpoints    crud.CRUDesque[uuid.UUID, stub.Endpoint]
	matchReports crud.CRUDesque[uuid.UUID, matching.Report]
	cdcService   CDCService
}

func NewService(cdcService CDCService) *Service {
	reportsCRUD := crud.New[uuid.UUID, matching.Report](matching.SortReport)
	endpointCRUD := crud.New[uuid.UUID, stub.Endpoint](stub.SortEndpoint)

	return &Service{
		endpoints:    endpointCRUD,
		matchReports: reportsCRUD,
		cdcService:   cdcService,
	}
}

func (m *Service) CheckEndpoints() ([]contract.Report, error) {
	var allReports []contract.Report

	endpoints, err := m.endpoints.GetAll()
	if err != nil {
		return nil, err
	}
	for _, endpoint := range endpoints {
		reports, err := m.cdcService.GetReports(endpoint)
		if err != nil {
			return nil, err
		}
		allReports = append(allReports, reports...)
	}
	return allReports, nil
}

func (m *Service) Reports() crud.CRUDesque[uuid.UUID, matching.Report] {
	return m.matchReports
}

func (m *Service) Endpoints() crud.CRUDesque[uuid.UUID, stub.Endpoint] {
	return m.endpoints
}

func (m *Service) CreateMatchReport(r *http.Request) (matching.Report, error) {
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
