package auth_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hreshchyshynt/chirpy/internal/auth"
)

func TestMakeJWT(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret"

	token, err := auth.MakeJWT(userID, secret, time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT() returned unexpected error: %v", err)
	}
	if token == "" {
		t.Fatal("MakeJWT() returned empty token")
	}
}

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret"

	token, err := auth.MakeJWT(userID, secret, time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT() returned unexpected error: %v", err)
	}

	gotID, err := auth.ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("ValidateJWT() returned unexpected error: %v", err)
	}
	if gotID != userID {
		t.Errorf("ValidateJWT() returned userID = %v, want %v", gotID, userID)
	}
}

func TestValidateJWTWrongSecret(t *testing.T) {
	userID := uuid.New()

	token, err := auth.MakeJWT(userID, "correct-secret", time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT() returned unexpected error: %v", err)
	}

	_, err = auth.ValidateJWT(token, "wrong-secret")
	if err == nil {
		t.Fatal("ValidateJWT() succeeded with wrong secret, expected error")
	}
}

func TestValidateJWTExpiredToken(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret"

	token, err := auth.MakeJWT(userID, secret, -time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT() returned unexpected error: %v", err)
	}

	_, err = auth.ValidateJWT(token, secret)
	if err == nil {
		t.Fatal("ValidateJWT() succeeded with expired token, expected error")
	}
}

func TestValidateJWTInvalidToken(t *testing.T) {
	_, err := auth.ValidateJWT("not-a-valid-jwt", "test-secret")
	if err == nil {
		t.Fatal("ValidateJWT() succeeded with invalid token string, expected error")
	}
}

func TestMakeJWTDifferentUsersProduceDifferentTokens(t *testing.T) {
	secret := "test-secret"
	user1 := uuid.New()
	user2 := uuid.New()

	token1, err := auth.MakeJWT(user1, secret, time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT() for user1 returned unexpected error: %v", err)
	}

	token2, err := auth.MakeJWT(user2, secret, time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT() for user2 returned unexpected error: %v", err)
	}

	if token1 == token2 {
		t.Error("MakeJWT() produced identical tokens for different users")
	}
}

func TestValidateJWTEmptyTokenString(t *testing.T) {
	_, err := auth.ValidateJWT("", "test-secret")
	if err == nil {
		t.Fatal("ValidateJWT() succeeded with empty token string, expected error")
	}
}

func TestValidateJWTCorrectUserIDReturned(t *testing.T) {
	secret := "test-secret"
	userIDs := []uuid.UUID{uuid.New(), uuid.New(), uuid.New()}

	for _, id := range userIDs {
		token, err := auth.MakeJWT(id, secret, time.Hour)
		if err != nil {
			t.Fatalf("MakeJWT() returned unexpected error: %v", err)
		}

		gotID, err := auth.ValidateJWT(token, secret)
		if err != nil {
			t.Fatalf("ValidateJWT() returned unexpected error: %v", err)
		}
		if gotID != id {
			t.Errorf("ValidateJWT() returned %v, want %v", gotID, id)
		}
	}
}

func TestValidateJWTDifferentSecretsSameUser(t *testing.T) {
	userID := uuid.New()

	token1, err := auth.MakeJWT(userID, "secret-one", time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT() returned unexpected error: %v", err)
	}

	token2, err := auth.MakeJWT(userID, "secret-two", time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT() returned unexpected error: %v", err)
	}

	if token1 == token2 {
		t.Error("MakeJWT() produced identical tokens with different secrets")
	}

	_, err = auth.ValidateJWT(token1, "secret-two")
	if err == nil {
		t.Error("ValidateJWT() succeeded validating token1 with secret-two, expected error")
	}

	_, err = auth.ValidateJWT(token2, "secret-one")
	if err == nil {
		t.Error("ValidateJWT() succeeded validating token2 with secret-one, expected error")
	}
}

func TestValidateJWTTokenNotYetExpired(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret"

	token, err := auth.MakeJWT(userID, secret, 2*time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT() returned unexpected error: %v", err)
	}

	gotID, err := auth.ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("ValidateJWT() returned unexpected error for valid token: %v", err)
	}
	if gotID != userID {
		t.Errorf("ValidateJWT() returned %v, want %v", gotID, userID)
	}
}
