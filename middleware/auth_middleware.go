package middleware

import (
	utils "golang_projects/utility"
	"net/http"
	"strings"
)

func JWTAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			//mak json response if
			utils.WriteJSONResponse(w, http.StatusUnauthorized, false, "Missing Authorization header", nil)
			return
		}

		// Extract token from "Bearer <token>"
		token := strings.TrimPrefix(authHeader, "Bearer ")
		_, err := utils.ValidateJWT(token)
		if err != nil {
			utils.WriteJSONResponse(w, http.StatusUnauthorized, false, "Invalid or expired token", nil)
			return
		}

		next.ServeHTTP(w, r)
	}

}
