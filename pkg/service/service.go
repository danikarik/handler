package service

import (
	"net/http"

	"golang.org/x/time/rate"
)

// Option used to modify service properties.
type Option func(*Service)

// WithRateBurst sets custom burst for rate limiter.
func WithRateBurst(n int) Option {
	return func(s *Service) {
		s.limiter = rate.NewLimiter(1, n)
	}
}

// Service holds application handlers.
type Service struct {
	limiter *rate.Limiter
	handler http.Handler
}

// New creates a new instance of `Service`.
func New(opts ...Option) *Service {
	srv := &Service{limiter: rate.NewLimiter(1, 100)}

	for _, opt := range opts {
		opt(srv)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", urlHandler)

	// Allow only HTTP POST.
	handler := allowedMethod("POST", mux)

	// Allow only JSON content type.
	handler = allowContentType("application/json", handler)

	// Rate limit.
	handler = rateLimit(srv.limiter, handler)

	srv.handler = handler
	return srv
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.handler.ServeHTTP(w, r)
}
