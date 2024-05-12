package middlewares

import (
	"net/http"
	"project01/src/auth"
	"project01/src/response"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := auth.ValidateToken(r); err != nil {
			response.ERROR(w, http.StatusUnauthorized, err)
			return
		}

		next(w, r)
	}
}
