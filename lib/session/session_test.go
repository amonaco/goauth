package session

import "testing"

func TestTokenIncludesUserAndCompanyIdentifiers(t *testing.T) {
	s := Session{
		ID:        "session-123",
		UserID:    42,
		CompanyID: 7,
	}

	if got, want := s.Token(), "42:7:session-123"; got != want {
		t.Fatalf("token mismatch: got %q, want %q", got, want)
	}
}
