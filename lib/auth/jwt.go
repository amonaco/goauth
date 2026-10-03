package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/amonaco/goauth/lib/config"
	"github.com/amonaco/goauth/lib/session"
)

type jwtClaims struct {
	UserID    uint32   `json:"user_id"`
	CompanyID uint32   `json:"company_id"`
	Roles     []string `json:"roles"`
	Exp       int64    `json:"exp"`
}

// GenerateJWT signs and returns a session JWT.
func GenerateJWT(sess session.Session) (string, error) {
	conf := config.Get()
	if conf.JWTSecret == "" || conf.JWTSecret == "change-me-in-production" {
		return "", errors.New("auth: jwt secret is unset")
	}

	claims := jwtClaims{
		UserID:    sess.UserID,
		CompanyID: sess.CompanyID,
		Roles:     sess.Roles,
		Exp:       time.Now().Add(time.Duration(conf.SessionTTL) * time.Second).Unix(),
	}

	payload, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
	encPayload := base64.RawURLEncoding.EncodeToString(payload)
	signature := signHS256(conf.JWTSecret, header+"."+encPayload)
	return header + "." + encPayload + "." + signature, nil
}

// ParseJWT validates a JWT and returns a session value.
func ParseJWT(token string) (session.Session, error) {
	conf := config.Get()
	if conf.JWTSecret == "" || conf.JWTSecret == "change-me-in-production" {
		return session.Session{}, errors.New("auth: jwt secret is unset")
	}

	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return session.Session{}, errors.New("auth: malformed jwt")
	}

	header, payload, signature := parts[0], parts[1], parts[2]
	expectedSignature := signHS256(conf.JWTSecret, header+"."+payload)
	if !hmac.Equal([]byte(signature), []byte(expectedSignature)) {
		return session.Session{}, errors.New("auth: jwt signature mismatch")
	}

	decodedPayload, err := base64.RawURLEncoding.DecodeString(payload)
	if err != nil {
		return session.Session{}, err
	}

	var claims jwtClaims
	if err := json.Unmarshal(decodedPayload, &claims); err != nil {
		return session.Session{}, err
	}

	if claims.Exp < time.Now().Unix() {
		return session.Session{}, errors.New("auth: jwt expired")
	}

	return session.Session{
		UserID:    claims.UserID,
		CompanyID: claims.CompanyID,
		Roles:     claims.Roles,
	}, nil
}

func signHS256(secret, payload string) string {
	h := hmac.New(sha256.New, []byte(secret))
	_, _ = h.Write([]byte(payload))
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}
