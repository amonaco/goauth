package auth

import (
	"context"
	"testing"

	"github.com/amonaco/goauth/lib/session"
)

func TestGetUserIDFromContext(t *testing.T) {
	ctx := context.WithValue(context.Background(), session.ContextKey("session"), session.Session{UserID: 13})

	userID, err := GetUserID(ctx)
	if err != nil {
		t.Fatalf("GetUserID returned unexpected error: %v", err)
	}
	if userID != 13 {
		t.Fatalf("user ID mismatch: got %d, want %d", userID, 13)
	}
}
