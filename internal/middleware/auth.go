package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/AleksandrCherepanov/tg-scheduler/internal/server"
)

func BasicAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if !strings.HasPrefix(req.URL.Path, "/api") {
			next.ServeHTTP(w, req)
			return
		}

		user, pass, ok := req.BasicAuth()
		if !ok {
			server.ResponseWithError(
				w,
				server.GetResponseError("Authorization header is not set", 401),
			)
			return
		}

		if user != os.Getenv("basicUser") || pass != os.Getenv("basicPass") {
			server.ResponseWithError(
				w,
				server.GetResponseError("Authorization header is not set", 401),
			)
			return
		}

		next.ServeHTTP(w, req)
	})
}
