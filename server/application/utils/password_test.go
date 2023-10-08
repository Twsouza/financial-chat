package utils

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	password := "password123"
	hashed, err := HashPassword(password)
	if err != nil {
		t.Errorf("HashPassword returned an error: %v", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	if err != nil {
		t.Errorf("Hashed password does not match original password: %v", err)
	}
}

func TestVerifyPassword(t *testing.T) {
	password := "password123"
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Errorf("Error generating hashed password: %v", err)
	}

	err = VerifyPassword(string(hashed), password)
	if err != nil {
		t.Errorf("VerifyPassword returned an error: %v", err)
	}
}
