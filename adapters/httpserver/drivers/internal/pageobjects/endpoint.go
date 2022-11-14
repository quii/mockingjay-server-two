package pageobjects

import (
	"fmt"
	"strconv"

	"github.com/go-rod/rod"
	"github.com/quii/mockingjay-server-two/domain/mockingjay"
)

func EndpointFromMarkup(el *rod.Element) (mockingjay.Endpoint, error) {
	getRequestField := func(field string) string {
		return el.MustElement(fmt.Sprintf(`[data-request-field=%s]`, field)).MustText()
	}
	getResponseField := func(field string) string {
		return el.MustElement(fmt.Sprintf(`[data-response-field=%s]`, field)).MustText()
	}
	statusText := getResponseField("status")
	statusCode, err := strconv.Atoi(statusText)
	if err != nil {
		return mockingjay.Endpoint{}, err
	}

	endpoint := mockingjay.Endpoint{
		Description: el.MustElement(`[data-field=description`).MustText(),
		Request: mockingjay.Request{
			Method:    getRequestField("method"),
			RegexPath: getRequestField("regexPath"),
			Path:      getRequestField("path"),
			Headers:   extractHeadersFromMarkup(el, `[data-request-field=headers] dl`),
			Body:      getRequestField("body"),
		},
		Response: mockingjay.Response{
			Status:  statusCode,
			Body:    getResponseField("body"),
			Headers: extractHeadersFromMarkup(el, `[data-response-field=headers] dl`),
		},
	}
	return endpoint, nil
}

func extractHeadersFromMarkup(el *rod.Element, selector string) mockingjay.Headers {
	var requestHeaders mockingjay.Headers
	dl := el.MustElement(selector)
	listItems := dl.MustElements("*")

	if len(listItems) > 0 {
		requestHeaders = make(mockingjay.Headers)
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
