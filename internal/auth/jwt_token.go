package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtToken struct {
	tokenSignature *[]byte
	expiration     time.Duration
}

func NewJwtToken(tokenSignature *[]byte, expiration time.Duration) *JwtToken {
	return &JwtToken{
		tokenSignature: tokenSignature,
		expiration:     expiration,
	}
}

func (jwtToken *JwtToken) CreateToken(username *string, token *string) error {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"authorized": true,
		"user":       *username,
		"exp":        time.Now().Add(jwtToken.expiration).Unix(),
		"iat":        time.Now(),
	})
	tokenResult, err := claims.SignedString(*jwtToken.tokenSignature)
	if err != nil {
		return err
	}
	*token = tokenResult
	return nil
}

func (jwtToken *JwtToken) VerifyToken(token *string) error {
	tokenResult, err := jwt.Parse(*token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signature method: %v", token.Header["alg"])
		}
		return *jwtToken.tokenSignature, nil
	})

	if err != nil {
		return err
	}

	if !tokenResult.Valid {
		return errors.New("the token is invalid")
	}

	claims, ok := tokenResult.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("an error occurred while accessing the token's claims")
	}

	expRaw, ok := claims["exp"]
	if !ok {
		return errors.New("expiration claim (exp) not found in token")
	}

	exp, ok := expRaw.(float64)
	if !ok {
		return errors.New("expiration claim (exp) is not a valid number")
	}

	expirationTime := time.Unix(int64(exp), 0)
	if expirationTime.Before(time.Now()) {
		return errors.New("the access token is expired")
	}

	return nil
}
