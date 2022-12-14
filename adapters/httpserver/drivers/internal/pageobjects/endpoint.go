package pageobjects

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/google/uuid"
	"github.com/quii/mockingjay-server-two/domain/mockingjay/stub"
)

func EndpointFromMarkup(el *rod.Element) (stub.Endpoint, error) {
	getRequestField := func(field string) string {
		return el.MustElement(fmt.Sprintf(`[data-request-field=%s]`, field)).MustText()
	}
	getResponseField := func(field string) string {
		return el.MustElement(fmt.Sprintf(`[data-response-field=%s]`, field)).MustText()
	}
	statusText := getResponseField("status")
	statusCode, err := strconv.Atoi(statusText)
	if err != nil {
		return stub.Endpoint{}, err
	}

	rawUUID := el.MustAttribute(`data-id`)
	if rawUUID == nil {
		return stub.Endpoint{}, errors.New("markup must be broken, no data-id attribute found")
	}

	id, err := uuid.Parse(*rawUUID)
	if err != nil {
		return stub.Endpoint{}, err
	}

	if id.String() == "" {
		return stub.Endpoint{}, errors.New("couldnt get an id for endpoint")
	}

	endpoint := stub.Endpoint{
		ID:          id,
		Description: el.MustElement(`[data-field=description`).MustText(),
		Request: stub.Request{
			Method:    getRequestField("method"),
			RegexPath: getRequestField("regexPath"),
			Path:      getRequestField("path"),
			Headers:   extractHeadersFromMarkup(el, `[data-request-field=headers] dl`),
			Body:      getRequestField("body"),
		},
		Response: stub.Response{
			Status:  statusCode,
			Body:    getResponseField("body"),
			Headers: extractHeadersFromMarkup(el, `[data-response-field=headers] dl`),
		},
	}
	return endpoint, nil
}

func InsertEndpoint(form *rod.Element, endpoint stub.Endpoint) error {
	fillTextWith := func(name string, value string) {
		form.MustElement(fmt.Sprintf(`*[name="%s"]`, name)).MustInput(value)
	}
	addHeaders := func(prefix string, headers stub.Headers) {
		for k, v := range headers {
			form.MustElement(fmt.Sprintf(`*[name="%s.header.name"]`, prefix)).MustInput(k)
			form.MustElement(fmt.Sprintf(`*[name="%s.header.values"]`, prefix)).MustInput(strings.Join(v, "; "))
		}
	}

	fillTextWith("description", endpoint.Description)
	fillTextWith("path", endpoint.Request.Path)
	fillTextWith("regexpath", endpoint.Request.RegexPath)
	fillTextWith("request.body", endpoint.Request.Body)
	form.MustElement(`*[name="method"]`).MustSelect(endpoint.Request.Method)
	addHeaders("request", endpoint.Request.Headers)
	fillTextWith("status", fmt.Sprintf("%d", endpoint.Response.Status))
	fillTextWith("response.body", endpoint.Response.Body)
	addHeaders("response", endpoint.Response.Headers)

	submitButton, err := form.Element(`#submit`)
	if err != nil {
		return err
	}
	if err := submitButton.Click(proto.InputMouseButtonLeft, 1); err != nil {
		return err
	}
	return nil
}

func extractHeadersFromMarkup(el *rod.Element, selector string) stub.Headers {
	var requestHeaders stub.Headers
	dl := el.MustElement(selector)
	listItems := dl.MustElements("*")

	if len(listItems) > 0 {
		requestHeaders = make(stub.Headers)
		currentKey := ""
		for _, item := range listItems {
			if item.String() == "<dt>" {
				currentKey = item.MustText()
				continue
			}
			requestHeaders[currentKey] = append(requestHeaders[currentKey], item.MustText())
		}
	}
	return requestHeaders
}
