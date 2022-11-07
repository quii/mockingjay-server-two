package httpserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/quii/mockingjay-server-two/domain/mockingjay"
	"github.com/quii/mockingjay-server-two/domain/mockingjay/matching"
)

/*
notes
- we can have mj listen on two ports, one for mj mgmt and the other for the stub server so we can configure without conflict. for now though, just the one
*/

type Driver struct {
	StubServerURL   string
	ConfigServerURL string
	Client          *http.Client
}

func (d Driver) Send(request mockingjay.Request) (mockingjay.Response, matching.Report, error) {
	var matchReport matching.Report

	req := request.ToHTTPRequest(d.StubServerURL)

	res, err := d.Client.Do(req)
	if err != nil {
		return mockingjay.Response{}, matchReport, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return mockingjay.Response{}, matchReport, err
	}

	_ = json.Unmarshal(body, &matchReport)

	return mockingjay.Response{
		Status:  res.StatusCode,
		Body:    string(body),
		Headers: mockingjay.Headers(res.Header),
	}, matchReport, nil
}

func (d Driver) Configure(es ...mockingjay.Endpoint) error {
	endpointJSON, err := json.Marshal(es)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, d.ConfigServerURL, bytes.NewReader(endpointJSON))
	if err != nil {
		return err
	}

	res, err := d.Client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusAccepted {
		return fmt.Errorf("got unexpected %d when trying to configure mj", res.StatusCode)
	}

	return nil
}
