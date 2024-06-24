package auth

import (
	"testing"
	"time"
)

func TestJwtToken(t *testing.T) {
	secret := []byte("secretKey")

	username := "testuser"

	jwtToken := NewJwtToken(&secret, time.Millisecond*500)

	var token string
	err := jwtToken.CreateToken(&username, &token)
	if err != nil {
		t.Fatalf("Error creating token: %v", err)
	}

	if token == "" {
		t.Fatal("Token was not created correctly")
	}

	// Wait to ensure the token has expired
	time.Sleep(time.Millisecond * 600)

	err = jwtToken.VerifyToken(&token)
	if err == nil {
		t.Fatal("Expired token was incorrectly verified as valid")
	}

	err = jwtToken.CreateToken(&username, &token)
	if err != nil {
		t.Fatalf("Error creating token: %v", err)
	}

	err = jwtToken.VerifyToken(&token)
	if err != nil {
		t.Fatalf("Error verifying token: %v", err)
	}
}
