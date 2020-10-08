package service_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/danikarik/handler/pkg/service"
)

func TestAllowedMethod(t *testing.T) {
	ts := httptest.NewServer(service.New())
	defer ts.Close()

	testCases := []struct {
		Name       string
		Method     string
		StatusCode int
	}{
		{
			Name:       "GET",
			Method:     "GET",
			StatusCode: http.StatusMethodNotAllowed,
		},
		{
			Name:       "PUT",
			Method:     "PUT",
			StatusCode: http.StatusMethodNotAllowed,
		},
		{
			Name:       "PATCH",
			Method:     "PATCH",
			StatusCode: http.StatusMethodNotAllowed,
		},
		{
			Name:       "DELETE",
			Method:     "DELETE",
			StatusCode: http.StatusMethodNotAllowed,
		},
		{
			Name:       "POST",
			Method:     "POST",
			StatusCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			req, err := http.NewRequest(tc.Method, ts.URL, strings.NewReader(`[]`))
			if err != nil {
				log.Fatalf("got error: %v", err)
			}
			req.Header.Set("Content-Type", "application/json")

			resp, err := ts.Client().Do(req)
			if err != nil {
				log.Fatalf("got error: %v", err)
			}

			if resp.StatusCode != tc.StatusCode {
				log.Fatalf("got: %v, expected: %v", resp.StatusCode, tc.StatusCode)
			}
		})
	}
}

func TestAllowedContentType(t *testing.T) {
	ts := httptest.NewServer(service.New())
	defer ts.Close()

	testCases := []struct {
		Name       string
		Method     string
		StatusCode int
	}{
		{
			Name:       "GET",
			Method:     "GET",
			StatusCode: http.StatusNotAcceptable,
		},
		{
			Name:       "POST",
			Method:     "POST",
			StatusCode: http.StatusNotAcceptable,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			req, err := http.NewRequest(tc.Method, ts.URL, strings.NewReader(`[]`))
			if err != nil {
				log.Fatalf("got error: %v", err)
			}

			resp, err := ts.Client().Do(req)
			if err != nil {
				log.Fatalf("got error: %v", err)
			}

			if resp.StatusCode != tc.StatusCode {
				log.Fatalf("got: %v, expected: %v", resp.StatusCode, tc.StatusCode)
			}
		})
	}
}

func TestRateLimit(t *testing.T) {
	testCases := []struct {
		Name       string
		Count      int
		Burst      int
		StatusCode int
	}{
		{
			Name:       "Exceeded",
			Count:      5,
			Burst:      3,
			StatusCode: http.StatusTooManyRequests,
		},
		{
			Name:       "OK",
			Count:      5,
			Burst:      100,
			StatusCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			ts := httptest.NewServer(service.New(service.WithRateBurst(tc.Burst)))
			defer ts.Close()

			statusCode := 0

			for i := 0; i < tc.Count; i++ {
				resp, err := ts.Client().Post(ts.URL, "application/json", strings.NewReader(`[]`))
				if err != nil {
					log.Fatalf("got error: %v", err)
				}
				statusCode = resp.StatusCode
			}

			if statusCode != tc.StatusCode {
				log.Fatalf("got: %v, expected: %v", statusCode, tc.StatusCode)
			}
		})
	}
}
