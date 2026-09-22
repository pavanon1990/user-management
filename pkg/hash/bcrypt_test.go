package hash_test

import (
	"testing"
	"user-management-api/pkg/hash"
)

func TestHashPassword(t *testing.T) {
	t.Run("hash is not same password", func(t *testing.T) {
		hashed, err := hash.HashPassword("secret123")
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		if hashed == "secret123" {
			t.Fatal("expected hashed password to differ from the password")
		}
	})

	t.Run("hashing the same password but different results", func(t *testing.T) {
		hashed1, err := hash.HashPassword("secret123")
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		hashed2, err := hash.HashPassword("secret123")
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		if hashed1 == hashed2 {
			t.Fatal("expected two hashes of the same password to differ (random salt)")
		}
	})
}

func TestComparePassword(t *testing.T) {
	t.Run("correct password matches", func(t *testing.T) {
		hashed, err := hash.HashPassword("secret123")
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}

		if err := hash.ComparePassword(hashed, "secret123"); err != nil {
			t.Fatalf("expected password to match, err: %v", err)
		}
	})

	t.Run("wrong password does not match", func(t *testing.T) {
		hashed, err := hash.HashPassword("secret123")
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}

		if err := hash.ComparePassword(hashed, "wrong-password"); err == nil {
			t.Fatal("expected an error for a wrong password")
		}
	})
}
