package app

import (
	"fmt"
	"net/http"
)

func TestMock(service MyService) string {
	return service.GetData()
}

// MyFunction uses an HTTP client to make a request.
func MyFunction(client HTTPClient) (string, error) {
	// Create an HTTP request
	req, err := http.NewRequest("GET", "https://example.com", nil)
	if err != nil {
		return "", err
	}

	// Use the provided HTTP client to perform the request
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Process the response (in this example, we just return the status code as a string)
	return fmt.Sprintf("Status code: %d", resp.StatusCode), nil
}
