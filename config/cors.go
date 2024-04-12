package config

import (
	"net/http"
	"strings"
)

func SetCORS(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		headers := w.Header()
		origin := r.Header.Get("Origin")
		allowed := []string{"GET", "POST"}

		if origin != "" {
			headers.Set("Access-Control-Allow-Origin", "*")
			if strings.HasPrefix(origin, "http") || strings.HasPrefix(origin, "https") {
				headers.Set("Access-Control-Allow-Origin", origin)
			}
			headers.Set("Access-Control-Allow-Headers", "*")
			headers.Set("Access-Control-Allow-Methods", strings.Join(allowed, ", "))
		}

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
