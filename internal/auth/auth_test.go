package auth

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	h1 := hashPassword("mypassword", "mysecret")
	h2 := hashPassword("mypassword", "mysecret")
	if h1 != h2 {
		t.Fatal("hash not deterministic")
	}
	if h1 == "" {
		t.Fatal("hash is empty")
	}
}

func TestTokenRoundtrip(t *testing.T) {
	secret := "test-auth-secret"
	token, err := CreateToken("alice", secret)
	if err != nil {
		t.Fatalf("CreateToken failed: %v", err)
	}
	if token == "" {
		t.Fatal("token is empty")
	}

	ok, userID, err := VerifyToken(token, secret, 120)
	if err != nil {
		t.Fatalf("VerifyToken failed: %v", err)
	}
	if !ok {
		t.Fatal("token verification failed")
	}
	if userID != "alice" {
		t.Fatalf("expected userID alice, got %s", userID)
	}
}
