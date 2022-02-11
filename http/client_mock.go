package client

import "net/http"

type ClientMock struct {
	DoFunc func(request *http.Request) (*http.Response, error)
}

func (c *ClientMock) Do(request *http.Request) (*http.Response, error) {
	return c.DoFunc(request)
}
