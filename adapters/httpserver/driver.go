package httpserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/quii/mockingjay-server-two/domain/endpoints"
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

func (d Driver) Do(request endpoints.Request) (endpoints.Response, error) {
	req := request.ToHTTPRequest(d.StubServerURL)
	res, err := d.Client.Do(req)
	if err != nil {
		return endpoints.Response{}, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return endpoints.Response{}, err
	}

	return endpoints.Response{
		Status:  res.StatusCode,
		Body:    string(body),
		Headers: endpoints.Headers(res.Header),
	}, nil
}

func (d Driver) Configure(endpoints endpoints.Endpoints) error {
	endpointJSON, err := json.Marshal(endpoints)
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
