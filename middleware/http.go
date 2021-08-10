package middleware

import (
	"net/http"

	"github.com/blushft/redtape"
)

// NewHTTPMiddleware returns an http handler that evaluates policy before returning child handler.
func NewHTTPMiddleware(e redtape.Enforcer, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		meta := make(map[string]interface{})
		for k := range r.Header {
			meta[k] = r.Header.Get(k)
		}

		req := redtape.NewRequest(
			redtape.RequestContext(r.Context(), meta),
			redtape.RequestResource(r.URL.Path),
			redtape.RequestAction(r.Method),
			redtape.RequestSubject(requestSubject(r)))

		if err := e.Enforce(req); err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		h.ServeHTTP(w, r)
	})
}

func requestSubject(r *http.Request) redtape.Subject {
	return redtape.NewSubject(r.RemoteAddr,
		redtape.SubjectName(r.UserAgent()),
		redtape.SubjectMeta(map[string]interface{}{
			"referer": r.Referer,
		}),
	)
}
