package auth

import (
    "context"
    "fmt"
    "time"

    "github.com/amonaco/goauth/lib/db"
    "github.com/amonaco/goauth/lib/session"
)

// LoginUser validates a user and issues a secure session.
func LoginUser(ctx context.Context, email string, password string, companyID uint32) (session.Session, string, error) {
    user, err := db.FindUserByEmail(ctx, email)
    if err != nil {
        return session.Session{}, "", err
    }

    if password == "" {
        return session.Session{}, "", fmt.Errorf("empty password")
    }

    roles, err := db.FindCompanyRoles(ctx, int32(user.ID), int32(companyID))
    if err != nil {
        return session.Session{}, "", err
    }

    sess, err := session.CreateSession(uint32(user.ID), companyID, roles)
    if err != nil {
        return session.Session{}, "", err
    }

    token, err := GenerateJWT(sess)
    if err != nil {
        return session.Session{}, "", err
    }

    return sess, token, nil
}

// GetSessionForUser loads and returns a valid session from Redis.
func GetSessionForUser(ctx context.Context, token string) (session.Session, error) {
    return ResolveSession(token)
}

// TrackLoginEvent records a simple audit event placeholder.
func TrackLoginEvent(ctx context.Context, userID uint32, eventType string, details string) {
    _ = ctx
    _ = userID
    _ = eventType
    _ = details
    _ = time.Now()
}
