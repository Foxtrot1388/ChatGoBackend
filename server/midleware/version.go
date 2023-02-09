package midleware

import (
	"ChatGo/pkg/logging"
	"github.com/kataras/versioning"
	"net/http"
)

func VersionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		version := r.Header.Get("API-VERSION")
		if version == "" {
			version = "1"
		}
		logger := logging.GetLogger()
		logger.Debugf("set API-VERSION %s", version)
		r = r.WithContext(versioning.WithVersion(r.Context(), version))
		next.ServeHTTP(w, r)
	})
}
