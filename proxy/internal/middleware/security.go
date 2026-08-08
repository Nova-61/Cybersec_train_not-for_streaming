package middleware

import (
	"net/http"
)

var forbiddenHeaders = []string{
	"Authorization",
	"X-Internal",
	"X-Secret",
}

func Security(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, header := range forbiddenHeaders {
			if r.Header.Get(header) != "" {
				http.Error(w, "Forbidden header: "+header, http.StatusForbidden)
				return
			}
		}

		if r.Method != "GET" && r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		next.ServeHTTP(w, r)
	})
}
