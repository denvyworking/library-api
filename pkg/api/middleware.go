package api

import (
	"context"
	"net/http"
	"time"
)

func (api *api) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		api.logger.Info("Received request", "method", r.Method, "path", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")

		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()

		r = r.WithContext(ctx)

		if !api.RightAuth(w, r) {
			return
		}
		next.ServeHTTP(w, r)
	})
}
