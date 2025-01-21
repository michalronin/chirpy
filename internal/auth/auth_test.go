package auth

import (
	"github.com/google/uuid"
	"testing"
	"time"
)

func TestMakeJWT(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret"
	t.Run("valid token", func(t *testing.T) {
		token, err := MakeJWT(userID, secret, 1*time.Second)
		if err != nil {
			t.Fatalf("failed to create token: %v", err)
		}
		valid, err := ValidateJWT(token, secret)
		if err != nil {
			t.Fatalf("failed to validate token: %v", err)
		}
		if valid != userID {
			t.Errorf("got user ID %v, want %v", valid, userID)
		}
	})
}

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret"

	t.Run("wrong secret", func(t *testing.T) {
		// First create a valid token
		token, err := MakeJWT(userID, secret, time.Second)
		if err != nil {
			t.Fatalf("failed to create token: %v", err)
		}

		// Try to validate with wrong secret
		_, err = ValidateJWT(token, "wrong-secret")
		if err == nil {
			t.Error("expected error with wrong secret, got nil")
		}
	})

	t.Run("expired token", func(t *testing.T) {
		// Create a token that expires very quickly
		token, err := MakeJWT(userID, secret, time.Millisecond)
		if err != nil {
			t.Fatalf("failed to create token: %v", err)
		}

		// Wait for token to expire
		time.Sleep(2 * time.Millisecond)

		// Try to validate expired token
		_, err = ValidateJWT(token, secret)
		if err == nil {
			t.Error("expected error with expired token, got nil")
		}
	})
}
