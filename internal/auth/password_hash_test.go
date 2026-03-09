package auth_test

import (
	"testing"

	"github.com/hreshchyshynt/chirpy/internal/auth"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		password string
	}{
		{
			password: "12345Kja@asdf.am",
		},
		{
			password: "asldkfjasdkjfcm,vm",
		},
		{
			password: "1238172348123",
		},
		{
			password: "@*#&@*&#@&#@)#()@()",
		},
	}
	t.Run("check hashing password match", func(t *testing.T) {
		for _, tt := range tests {
			got, gotErr := auth.HashPassword(tt.password)
			if gotErr != nil {
				t.Errorf("Error received during hashing password: %v", gotErr)
			}
			match, err := auth.CheckPasswordHash(tt.password, got)
			if err != nil {
				t.Errorf("Error received during checking password: %v", err)
			}
			if !match {
				t.Errorf("HashPassword() = %v, match: %v", got, match)
			}

		}

	})

}

func TestCheckPasswordHashWrongPassword(t *testing.T) {
	hash, err := auth.HashPassword("correct-password")
	if err != nil {
		t.Fatalf("HashPassword() returned unexpected error: %v", err)
	}

	match, err := auth.CheckPasswordHash("wrong-password", hash)
	if err != nil {
		t.Fatalf("CheckPasswordHash() returned unexpected error: %v", err)
	}
	if match {
		t.Error("CheckPasswordHash() matched with wrong password, expected no match")
	}
}

func TestHashPasswordProducesDifferentHashes(t *testing.T) {
	password := "same-password"

	hash1, err := auth.HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() returned unexpected error: %v", err)
	}

	hash2, err := auth.HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() returned unexpected error: %v", err)
	}

	if hash1 == hash2 {
		t.Error("HashPassword() produced identical hashes for same password, expected unique salts")
	}
}

func TestHashPasswordNotPlaintext(t *testing.T) {
	password := "my-secret-password"

	hash, err := auth.HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() returned unexpected error: %v", err)
	}

	if hash == password {
		t.Error("HashPassword() returned the plaintext password")
	}
}

func TestCheckPasswordHashInvalidHash(t *testing.T) {
	_, err := auth.CheckPasswordHash("password", "not-a-valid-hash")
	if err == nil {
		t.Error("CheckPasswordHash() succeeded with invalid hash, expected error")
	}
}
