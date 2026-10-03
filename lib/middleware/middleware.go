package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/amonaco/goauth/lib/auth"
	"github.com/amonaco/goauth/lib/session"
)

// Middleware enforces a valid session cookie or auth-token header.
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var token string

		cookie, err := r.Cookie(auth.TokenCookieName)
		if err == nil {
			token = cookie.Value
		} else {
			token = r.Header.Get("auth-token")
			if token == "" {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
		}

		sess, err := session.GetSession(token)
		if err != nil {
			log.Println("auth middleware: invalid session:", err)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		ctx = context.WithValue(ctx, session.ContextKey("session"), sess)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
