package app

import (
	"net/http"

	"github.com/stretchr/testify/mock"
)

type Metrics struct {
	Name string
}

func (m *Metrics) GetData() string {
	return "mock"
}

// RealHTTPClient is the real implementation of HTTPClient.
type RealHTTPClient struct{}

// Do performs an HTTP request and returns the response.
func (c *RealHTTPClient) Do(req *http.Request) (*http.Response, error) {
	client := http.Client{}
	return client.Do(req)
}

// MockHTTPClient is a mock implementation of HTTPClient.
type MockHTTPClient struct {
	mock.Mock
}

// Do mocks the HTTP request.
func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}
