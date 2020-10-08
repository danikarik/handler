package service

import (
	"net/http"

	"golang.org/x/time/rate"
)

func allowedMethod(method string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			httpError(w, http.StatusMethodNotAllowed)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func allowContentType(contentType string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != contentType {
			httpError(w, http.StatusNotAcceptable)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func rateLimit(limiter *rate.Limiter, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			httpError(w, http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
