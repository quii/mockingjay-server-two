package mockingjay

import (
	"github.com/google/uuid"
	"github.com/quii/mockingjay-server-two/domain/collections"
	"github.com/quii/mockingjay-server-two/domain/mockingjay/contract"
	"github.com/quii/mockingjay-server-two/domain/mockingjay/matching"
	"github.com/quii/mockingjay-server-two/domain/mockingjay/stub"
)

type Driver struct {
	service *Service
}

func NewDriver(service *Service) *Driver {
	return &Driver{service: service}
}

func (d Driver) Send(request stub.Request) (matching.Report, error) {
	report, err := d.service.CreateMatchReport(request)
	return report, err
}

func (d Driver) CheckEndpoints() ([]contract.Report, error) {
	return d.service.CheckEndpoints()
}

func (d Driver) GetReports() ([]matching.Report, error) {
	return d.service.Reports().GetAll()
}

func (d Driver) DeleteReports() error {
	reports, err := d.service.Reports().GetAll()
	if err != nil {
		return err
	}
	return collections.ForAll(reports, func(r matching.Report) error {
		return d.service.Reports().Delete(r.ID)
	})
}

func (d Driver) AddEndpoints(endpoints ...stub.Endpoint) error {
	return collections.ForAll(endpoints, func(e stub.Endpoint) error {
		//todo: Maybe service shouldn't expose its crud, and instead have a method that encpasulates compiling, so the rest of the system doesn't have to know
		if err := e.Compile(); err != nil {
			return err
		}
		return d.service.Endpoints().Create(e.ID, e)
	})
}

func (d Driver) GetEndpoints() (stub.Endpoints, error) {
	return d.service.Endpoints().GetAll()
}

func (d Driver) DeleteEndpoint(uuid uuid.UUID) error {
	return d.service.Endpoints().Delete(uuid)
}

func (d Driver) DeleteEndpoints() error {
	endpoints, err := d.service.Endpoints().GetAll()
	if err != nil {
		return err
	}
	return collections.ForAll(endpoints, func(e stub.Endpoint) error {
		return d.service.Reports().Delete(e.ID)
	})
}
