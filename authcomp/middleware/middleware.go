package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/melsonic/skyvault/auth/db"
	jwtauth "github.com/melsonic/skyvault/auth/jwt"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("Authorization")
		tokenString := strings.Split(authorizationHeader, " ")[1]

		if authorizationHeader == "" || strings.HasPrefix(authorizationHeader, "Bearer ") || strings.Count(tokenString, ".") != 2 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid authorization header"))
			return
		}

		user := jwtauth.GetUserIdentityFromAccessToken(tokenString)

		if user == nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("invalid jwt token"))
			return
		}

		// fill other user details
		err := db.GetUserData(user)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("error fetching user data"))
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), "user", user))

		next(w, r)
	}
}
