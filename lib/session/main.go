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

// Session holds permission and identity informations about a user
type Session struct {
	ID        string
	Roles     []string
	UserID    uint32
	CompanyID uint32
}

// ContextKey is used as the key type to store
// a session in the request context
type ContextKey string

// Token returns a session token used to retreive the session from redis
func (s Session) Token() string {

	// The session ID follows the format:
	// user_id:company_id:session
	return fmt.Sprintf("%v:%v:%v", s.UserID, s.CompanyID, s.ID)
}

// GetSession fetches a session by token from redis
func GetSession(token string) (Session, error) {
	var session Session

	data, err := cache.Get(makeSessionKey(token))
	if err != nil {
		return session, err
	}

	err = json.Unmarshal([]byte(data), &session)
	if err != nil {
		return session, err
	}

	return session, nil
}

// CreateSession creates a new session and stores it in redis
func CreateSession(userID uint32, companyID uint32, roles []string) (Session, error) {
	id, err := generateSessionID()
	if err != nil {
		// Something is very wrong: OS could not generate a random number
		return Session{}, err
	}

	session := Session{
		ID:        id,
		UserID:    userID,
		Roles:     roles,
		CompanyID: companyID,
	}

	data, err := json.Marshal(session)
	if err != nil {
		return Session{}, err
	}

	// Store session in redis with a TTL of 24 hours
	err = cache.Set(makeSessionKey(session.ID), string(data), sessionTTL)
	if err != nil {
		return Session{}, err
	}

	// Add user tokens to a list to allow easy expiring sessions
	err = cache.PushExpire(makeUserKey(userID, companyID), string(id), sessionTTL)
	if err != nil {
		return Session{}, err
	}

	return session, nil
}

// DeleteSession removes a session from redis
func DeleteSession(token string) error {
	var session Session

	// GetDel the session and data
	data, err := cache.GetDel(makeSessionKey(token))
	if err != nil {
		return err
	}

	if len(data) == 0 {
		return nil
	}

	err = json.Unmarshal([]byte(data), &session)
	if err != nil {
		return err
	}

	// Remove this specific token from user session list
	err = cache.LRem(makeUserKey(session.UserID, session.CompanyID), token)
	if err != nil {
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
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}

// Generates a token used for signup and password recovery
func GenerateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}
