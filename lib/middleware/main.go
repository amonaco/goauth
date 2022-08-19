package auth

import (
	"context"
	"log"
	"net/http"
)

// Main middleware function
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var token string

		// Try cookie auth first
		cookie, err := r.Cookie(TokenCookieName)
		if err != nil {
			log.Println("Session cookie not present! Fallback to auth-token")

			// Try auth-token header
			token = r.Header.Get("auth-token")
			if token == "" {
				http.Error(w, http.StatusText(401), 401)
				return
			}
		} else {
			token = cookie.Value
		}

		// Check session token exists
		session, err := GetSession(token)
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(401), 401)
			return
		}

		ctx = context.WithValue(ctx, ContextKey("session"), session)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
