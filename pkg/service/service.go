package service

import (
	"net/http"
)

// Service holds application handlers.
type Service struct{ handler http.Handler }

// New creates a new instance of `Service`.
func New() *Service {
	mux := http.NewServeMux()
	mux.HandleFunc("/", urlHandler)

	// Allow only HTTP POST.
	handler := allowedMethod("POST", mux)

	// Allow only JSON content type.
	handler = allowContentType("application/json", handler)

	return &Service{handler}
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.handler.ServeHTTP(w, r)
}
