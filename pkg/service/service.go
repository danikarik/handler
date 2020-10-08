package service

import (
	"net/http"
	"time"
)

// Service holds application handlers.
type Service struct{ srv *http.Server }

// New creates a new instance of `Service`.
func New(addr string) *Service {
	mux := http.NewServeMux()
	mux.HandleFunc("/", urlHandler)

	// Allow only HTTP POST.
	handler := allowedMethod("POST", mux)

	// Allow only JSON content type.
	handler = allowContentType("application/json", handler)

	return &Service{
		srv: &http.Server{
			Addr:         addr,
			Handler:      handler,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	}
}

// Start runs http server.
func (s *Service) Start() error {
	return s.srv.ListenAndServe()
}
