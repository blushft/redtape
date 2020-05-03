package middleware

import (
	"net/http"

	"github.com/blushft/redtape"
)

// NewHTTPMiddleware returns an http handler that evaluates policy before returning child handler
func NewHTTPMiddleware(e redtape.Enforcer, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := redtape.NewRequestWithContext(r.Context(), r.URL.Path, r.Method, "", "", requestMetadata(r))

		if err := e.Enforce(req); err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		h.ServeHTTP(w, r)
	})
}

func requestMetadata(r *http.Request) map[string]interface{} {

	return map[string]interface{}{
		"referer":    r.Referer,
		"cookies":    r.Cookies(),
		"user_agent": r.UserAgent(),
		"url":        r.URL.String,
		"headers":    r.Header,
	}
}
