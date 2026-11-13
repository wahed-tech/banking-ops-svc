package api

import (
	"net/http"
	"slices"
	"strings"
)

// preflightMaxAge - seconds a browser may reuse a preflight answer before asking again
const preflightMaxAge = "3600"

// withCORS - lets only the ops portal's own origins call the API from a browser
// Preflight (OPTIONS) is answered here, before auth, because browsers send it without the token
func withCORS(allowedOrigins []string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		allowed := origin != "" && slices.Contains(allowedOrigins, origin)

		if allowed {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			w.Header().Add("Vary", "Origin")
		}

		if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
			if !allowed {
				w.WriteHeader(http.StatusForbidden)
				return
			}
			w.Header().Set("Access-Control-Max-Age", preflightMaxAge)
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// ParseOrigins - splits a comma-separated origin list from config, so it can be set from a single env var
func ParseOrigins(raw string) []string {
	var origins []string
	for o := range strings.SplitSeq(raw, ",") {
		if o = strings.TrimSpace(o); o != "" {
			origins = append(origins, o)
		}
	}

	return origins
}
