package middlewares

import (
	"net/http"

	"go.uber.org/zap"
)

func RecoverPanic(logger *zap.SugaredLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					w.Header().Set("Connection", "close")
					logger.Error(err)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
