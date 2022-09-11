package middleware

import (
	"go/token"
	"net/http"
	"strings"

	"github.com/Euler-B/API-REST_Go/models"
	"github.com/Euler-B/API-REST_Go/server"
	"github.com/golang-jwt/jwt/v4"
)

var (
	NO_AUTH_NEEDED = []string {
		"login",
		"signup",
	}
)

func shouldCheckToken(route string) bool {
	for _, p := range NO_AUTH_NEEDED{
		if strings.Contains(route, p) {
			return false
		}
	}
	return true
}

func CheckAuthMiddleware(s server.Server) func(h http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !shouldCheckToken(r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}
			tokenString := strings.TrimSpace(r.Header.Get("Autorization"))
			_, err := jwt.ParseWithClaims(tokenString, models.AppClaims{}, 
			func(token *jwt.Token) (interface{}, error {
				return []byte(s.Config().JWTSecret), nil
			})x
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}