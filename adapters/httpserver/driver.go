package httpserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/quii/mockingjay-server-two/domain/config"
)

/*
notes
- we can have mj listen on two ports, one for config mgmt and the other for the stub server so we can configure without conflict. for now though, just the one
*/

type Driver struct {
	StubServerURL string
	Client        *http.Client
}

func (d Driver) Do(request config.Request) (config.Response, error) {
	req, err := http.NewRequest(request.Method, d.StubServerURL+request.Path, nil)
	if err != nil {
		return config.Response{}, err
	}

	res, err := d.Client.Do(req)
	if err != nil {
		return config.Response{}, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return config.Response{}, err
	}

	return config.Response{
		Status: res.StatusCode,
		Body:   string(body),
	}, nil
}

func (d Driver) Configure(endpoints config.Endpoints) error {
	endpointJSON, err := json.Marshal(endpoints)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, d.StubServerURL+"/configure", bytes.NewReader(endpointJSON))
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
