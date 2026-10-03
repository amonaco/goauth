package auth

import (
	"net/http"
	"strings"

	"github.com/amonaco/goauth/lib/session"
)

// SetSessionCookie writes a secure, HttpOnly session cookie.
func SetSessionCookie(w http.ResponseWriter, token string, ttl int, secure bool) {
	cookie := &http.Cookie{
		Name:     TokenCookieName,
		Value:    token,
		Path:     "/",
		MaxAge:   ttl,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, cookie)
}

// ResolveSession accepts either a session token stored in Redis or a signed JWT.
func ResolveSession(token string) (session.Session, error) {
	if strings.Count(token, ".") == 2 {
		if sess, err := ParseJWT(token); err == nil {
			return sess, nil
		}
	}
	return session.GetSession(token)
}
