package mockingjay

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/quii/mockingjay-server-two/domain/crud"
	"github.com/quii/mockingjay-server-two/domain/mockingjay/contract"
	http2 "github.com/quii/mockingjay-server-two/domain/mockingjay/http"
	"github.com/quii/mockingjay-server-two/domain/mockingjay/matching"
)

type CDCService interface {
	GetReports(endpoint http2.Endpoint) ([]contract.Report, error)
}

type Service struct {
	endpoints    crud.CRUDesque[uuid.UUID, http2.Endpoint]
	matchReports crud.CRUDesque[uuid.UUID, matching.Report]
	cdcService   CDCService
}

func NewService(endpoints http2.Endpoints, cdcService CDCService) (*Service, error) {
	reportsCRUD := crud.New[uuid.UUID, matching.Report](matching.SortReport)
	endpointCRUD := crud.New[uuid.UUID, http2.Endpoint](http2.SortEndpoint)

	for _, endpoint := range endpoints {
		if err := endpointCRUD.Create(endpoint.ID, endpoint); err != nil {
			return nil, err
		}
	}
	return &Service{
		endpoints:    endpointCRUD,
		matchReports: reportsCRUD,
		cdcService:   cdcService,
	}, nil
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

func (m *Service) Endpoints() crud.CRUDesque[uuid.UUID, http2.Endpoint] {
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
