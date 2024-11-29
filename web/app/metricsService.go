package app

import (
	"net/http"
)

type MyService interface {
	GetData() string
}
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}
