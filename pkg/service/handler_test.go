package service_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/danikarik/handler/pkg/service"
)

func TestUrlHandlerRequest(t *testing.T) {
	ts := httptest.NewServer(service.New())
	defer ts.Close()

	testCases := []struct {
		Name       string
		Body       string
		StatusCode int
	}{
		{
			Name:       "OK",
			Body:       `["https://kaspi.kz"]`,
			StatusCode: http.StatusOK,
		},
		{
			Name:       "Empty",
			Body:       `[]`,
			StatusCode: http.StatusOK,
		},
		{
			Name: "AboveLimit",
			Body: `[
						"https://google.com", "https://apple.com",
						"https://google.com", "https://apple.com",
						"https://google.com", "https://apple.com",
						"https://google.com", "https://apple.com",
						"https://google.com", "https://apple.com",
						"https://google.com", "https://apple.com",
						"https://google.com", "https://apple.com",
						"https://google.com", "https://apple.com",
						"https://google.com", "https://apple.com",
						"https://google.com", "https://apple.com",
						"https://google.com", "https://apple.com",
						"https://google.com", "https://apple.com",
					]`,
			StatusCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			resp, err := ts.Client().Post(ts.URL, "application/json", strings.NewReader(tc.Body))
			if err != nil {
				log.Fatalf("got error: %v", err)
			}

			if resp.StatusCode != tc.StatusCode {
				log.Fatalf("got: %v, expected: %v", resp.StatusCode, tc.StatusCode)
			}
		})
	}
}
