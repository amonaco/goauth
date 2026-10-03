package session

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/amonaco/goauth/lib/cache"
)

const sessionTTL = 86400
const TokenExpiry = 86400

// Session holds permission and identity information about a user.
type Session struct {
	ID        string
	Roles     []string
	UserID    uint32
	CompanyID uint32
}

// ContextKey is used as the key type to store a session in the request context.
type ContextKey string

// Token returns a session token used to retrieve the session from Redis.
func (s Session) Token() string {
	return fmt.Sprintf("%v:%v:%v", s.UserID, s.CompanyID, s.ID)
}

// GetSession fetches a session by token from Redis.
func GetSession(token string) (Session, error) {
	var sessionValue Session

	data, err := cache.Get(makeSessionKey(token))
	if err != nil {
		return sessionValue, err
	}

	if err := json.Unmarshal([]byte(data), &sessionValue); err != nil {
		return sessionValue, err
	}

	return sessionValue, nil
}

// CreateSession creates a new session and stores it in Redis.
func CreateSession(userID uint32, companyID uint32, roles []string) (Session, error) {
	id, err := generateSessionID()
	if err != nil {
		return Session{}, err
	}

	sessionValue := Session{
		ID:        id,
		UserID:    userID,
		Roles:     roles,
		CompanyID: companyID,
	}

	data, err := json.Marshal(sessionValue)
	if err != nil {
		return Session{}, err
	}

	key := makeSessionKey(sessionValue.Token())
	if err := cache.Set(key, string(data), sessionTTL); err != nil {
		return Session{}, err
	}

	if err := cache.PushExpire(makeUserKey(userID, companyID), sessionValue.Token(), sessionTTL); err != nil {
		return Session{}, err
	}

	return sessionValue, nil
}

// DeleteSession removes a session from Redis.
func DeleteSession(token string) error {
	var sessionValue Session

	data, err := cache.GetDel(makeSessionKey(token))
	if err != nil {
		return err
	}

	if len(data) == 0 {
		return nil
	}

	if err := json.Unmarshal([]byte(data), &sessionValue); err != nil {
		return err
	}

	if err := cache.LRem(makeUserKey(sessionValue.UserID, sessionValue.CompanyID), token); err != nil {
		return err
	}

	return nil
}

func makeUserKey(userID uint32, companyID uint32) string {
	return fmt.Sprintf("user:%d:%d", userID, companyID)
}

func makeSessionKey(token string) string {
	return "session:" + token
}

func generateSessionID() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}

// GenerateToken creates a cryptographically random token for signup and password recovery flows.
func GenerateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}
