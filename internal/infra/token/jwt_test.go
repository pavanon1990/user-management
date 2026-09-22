package token_test

import (
	"testing"
	"time"

	"user-management-api/internal/infra/token"
)

func TestJWTProvider_ValidateToken(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		provider := token.NewJWTProvider("secret-key", time.Hour)

		tokenStr, err := provider.GenerateToken("uuid-1")
		if err != nil {
			t.Fatalf("unexpected err generating token: %v", err)
		}

		userID, err := provider.ValidateToken(tokenStr)
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		if userID != "uuid-1" {
			t.Fatalf("expected userID uuid-1")
		}
	})

	t.Run("expired token", func(t *testing.T) {
		provider := token.NewJWTProvider("secret-key", -time.Hour)

		tokenStr, err := provider.GenerateToken("uuid-1")
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}

		_, err = provider.ValidateToken(tokenStr)
		if err == nil {
			t.Fatal("expected error for expired token")
		}
	})
}
