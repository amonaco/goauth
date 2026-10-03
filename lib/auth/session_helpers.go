package auth

import (
	"encoding/json"
	"net/http"

	"github.com/amonaco/goauth/lib/config"
	"github.com/amonaco/goauth/lib/session"
)

// IssueSession creates a Redis-backed session and returns a signed JWT token.
func IssueSession(w http.ResponseWriter, r *http.Request, userID, companyID uint32, roles []string) (session.Session, string, error) {
	sess, err := session.CreateSession(userID, companyID, roles)
	if err != nil {
		return session.Session{}, "", err
	}

	token, err := GenerateJWT(sess)
	if err != nil {
		return session.Session{}, "", err
	}

	conf := config.Get()
	SetSessionCookie(w, token, conf.SessionTTL, conf.CookieSecure)
	return sess, token, nil
}

// RefreshSession rotates a session token while preserving identity and roles.
func RefreshSession(currentToken string) (string, error) {
	sess, err := ResolveSession(currentToken)
	if err != nil {
		return "", err
	}

	return GenerateJWT(sess)
}

// Logout clears the caller's session from Redis and expires the cookie when present.
func Logout(w http.ResponseWriter, r *http.Request) error {
	token := r.Header.Get("auth-token")
	if token == "" {
		if cookie, err := r.Cookie(TokenCookieName); err == nil {
			token = cookie.Value
		}
	}

	if token != "" {
		if _, err := ParseJWT(token); err == nil {
			SetSessionCookie(w, "", -1, config.Get().CookieSecure)
			return nil
		}
		if err := session.DeleteSession(token); err != nil {
			return err
		}
	}

	SetSessionCookie(w, "", -1, config.Get().CookieSecure)
	return nil
}

// JSONResponse writes a JSON payload to a response.
func JSONResponse(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
