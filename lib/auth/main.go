package auth

import (
	"context"
	"errors"
	"log"

	"github.com/amonaco/goauth/lib/session"
)

const TokenCookieName = "goauth_session"

type Role map[string][]uint32

type Auth struct {
	ID     uint32 `json:"id"`
	UserID uint32 `json:"user_id"`
	Roles  Role   `json:"permission"`
}

// Authorizes a given session.
func Authorize(ctx context.Context, permissions ...string) error {
	if !isAuthorized(ctx, permissions) {
		return errors.New("auth: unauthorized")
	}

	return nil
}

// TODO: Change isAuthorized to isValid.
func isAuthorized(ctx context.Context, permissions []string) bool {
	sessionValue, ok := getSessionFromContext(ctx)
	if !ok {
		log.Println("auth: no session in context")
		return false
	}

	for _, permission := range permissions {
		for _, role := range sessionValue.Roles {
			if role == permission || role == "superadmin" {
				return true
			}
		}
	}

	return false
}

// GetUserID fetches the user ID from the session.
func GetUserID(ctx context.Context) (uint32, error) {
	sessionValue, ok := getSessionFromContext(ctx)
	if !ok {
		return 0, errors.New("auth: invalid user id")
	}

	return sessionValue.UserID, nil
}

// IsSuperAdmin checks if the current session is a superadmin.
func IsSuperAdmin(ctx context.Context) bool {
	sessionValue, ok := getSessionFromContext(ctx)
	if !ok {
		return false
	}

	for _, role := range sessionValue.Roles {
		if role == "superadmin" {
			return true
		}
	}

	return false
}

func getSessionFromContext(ctx context.Context) (session.Session, bool) {
	if ctx == nil {
		return session.Session{}, false
	}

	value, ok := ctx.Value(session.ContextKey("session")).(session.Session)
	return value, ok
}
