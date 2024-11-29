package app

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMetricsService(t *testing.T) {
	metric := Metrics{Name: "name"}
	result := metric.GetData()
	expected := "mock"
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}
func TestMyFunction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockMyService(ctrl)
	mockService.EXPECT().GetData().Return("mocked")

	// Now use the mockService in your code that depends on MyService
	result := TestMock(mockService)
	expected := "mocked"
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	} else {
		t.Log(">>>>>>>>>.")
	}
}

func TestMyFunctionWithMock(t *testing.T) {
	mockClient := new(MockHTTPClient)
	mockClient.On("Do", mock.Anything).Return(
		&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader("status")),
		},
		nil,
	)
	result, err := MyFunction(mockClient)

	assert.NoError(t, err)
	assert.Equal(t, "Status code: 200", result)

	mockClient.AssertExpectations(t)
}
