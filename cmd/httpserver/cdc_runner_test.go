package main_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"cuelang.org/go/cue/cuecontext"
	"github.com/alecthomas/assert/v2"
	"github.com/google/uuid"
	"github.com/quii/mockingjay-server-two/adapters"
	"github.com/quii/mockingjay-server-two/domain/mockingjay/contract"
	"github.com/quii/mockingjay-server-two/domain/mockingjay/matching"
	"github.com/quii/mockingjay-server-two/domain/mockingjay/stub"
	"github.com/quii/mockingjay-server-two/specifications"
)

func TestMockingjayCDCRunner(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	driver := &CDCDriver{}

	specifications.MockingjayConsumerDrivenContractSpec(t, driver, driver, specRoot)

	t.Cleanup(func() {
		assert.NoError(t, os.RemoveAll(driver.EndpointDirName))
	})
}

type CDCDriver struct {
	EndpointDirName string
	endpoints       []stub.Endpoint
}

func (C *CDCDriver) CheckEndpoints() ([]contract.Report, error) {
	stdout, err := adapters.RunMockingjayCLI(C.EndpointDirName)
	if err != nil {
		return nil, err
	}

	var reports []contract.Report
	err = json.Unmarshal([]byte(stdout), &reports)
	return reports, err
}

func (C *CDCDriver) AddEndpoints(endpoints ...stub.Endpoint) error {
	C.endpoints = endpoints
	endpointsCue := stub.EndpointsCue{Endpoints: endpoints}
	c := cuecontext.New()
	values := c.Encode(endpointsCue)

	dname, err := os.MkdirTemp("", "mj-blackbox-test")
	if err != nil {
		return err
	}
	C.EndpointDirName = dname

	fname := filepath.Join(dname, "endooints.cue")
	f, err := os.OpenFile(fname, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0700)
	if err != nil {
		return err
	}
	_, _ = fmt.Fprint(f, values)
	return nil
}

func (C *CDCDriver) GetEndpoints() (stub.Endpoints, error) {
	//todo: really, this should probably load the temp file we made to check it got written ok, but cba for now
	return C.endpoints, nil
}

func (C *CDCDriver) DeleteEndpoint(_ uuid.UUID) error {
	return nil
}

func (C *CDCDriver) GetReports() ([]matching.Report, error) {
	panic("implement me")
}

func (C *CDCDriver) DeleteEndpoints() error {
	return nil
}

func (C *CDCDriver) DeleteReports() error {
	return nil
}
