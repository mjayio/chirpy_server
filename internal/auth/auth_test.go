package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestHashedPassword(t *testing.T) {
	password := "password123"
	hashedPassword, err := HashedPassword(password)
	if err != nil {
		t.Fatalf("HashedPassword failed: %v", err)
	}
	if hashedPassword == password {
		t.Error("HashedPassword should not return the same password")
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "password123"
	hashedPassword, err := HashedPassword(password)
	if err != nil {
		t.Fatalf("HashedPassword failed: %v", err)
	}

	err = CheckPasswordHash(password, hashedPassword)
	if err != nil {
		t.Errorf("CheckPasswordHash failed: %v", err)
	}

	wrongPassword := "wrongpassword"
	err = CheckPasswordHash(wrongPassword, hashedPassword)
	if err == nil {
		t.Error("CheckPasswordHash should have failed for wrong password")
	}
}

func TestMakeJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "testsecret"
	expiresIn := 1 * time.Hour

	jwtToken, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	if jwtToken == "" {
		t.Error("MakeJWT should not return an empty string")
	}
}

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "testsecret"
	expiresIn := 1 * time.Hour

	jwtToken, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	validatedUserID, err := ValidateJWT(jwtToken, tokenSecret)
	if err != nil {
		t.Fatalf("ValidateJWT failed: %v", err)
	}

	if validatedUserID != userID {
		t.Errorf("ValidateJWT returned wrong user ID: got %v, want %v", validatedUserID, userID)
	}

	// Test with expired token
	expiredToken, err := MakeJWT(userID, tokenSecret, -1*time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	_, err = ValidateJWT(expiredToken, tokenSecret)
	if err == nil {
		t.Error("ValidateJWT should have failed for expired token")
	}

	// Test with wrong secret
	_, err = ValidateJWT(jwtToken, "wrongsecret")
	if err == nil {
		t.Error("ValidateJWT should have failed for wrong secret")
	}

	// Test with invalid token
	_, err = ValidateJWT("invalidtoken", tokenSecret)
	if err == nil {
		t.Error("ValidateJWT should have failed for invalid token")
	}
}
