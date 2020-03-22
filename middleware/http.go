package middleware

import (
	"net/http"

	"github.com/blushft/redtape"
)

func NewHTTPMiddlware(e redtape.Enforcer, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := &redtape.Request{
			Action:   r.Method,
			Resource: r.URL.Path,
			Context: redtape.RequestContext{
				Context:  r.Context(),
				Metadata: requestMetadata(r),
			},
		}

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
