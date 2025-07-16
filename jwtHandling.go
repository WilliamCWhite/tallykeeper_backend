package main

import (
	"fmt"
	"os"
	"time"
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

// Claims I want in my jwt
type TallyJwtClaims struct {
	UserID               string `json:"user_id"`
	jwt.RegisteredClaims        // Embed standard claims in struct
}

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func GenerateJWT(userID string) (string, error) {
	claims := TallyJwtClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)), // HOW LONG TOKEN LASTS
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "tallykeeper_backend",
			Subject:   userID,
			Audience:  []string{"tallykeeper_frontend"},
		},
	}

	// Create token with claims and signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret key
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

func VerifyJWT(tokenString string) (*TallyJwtClaims, error) {
	claims := &TallyJwtClaims{}

	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		// ensure expected method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// send secret to ParseWithClaims function above
		return jwtSecret, nil
	})

	if !token.Valid {
		if err != nil {
			if errors.Is(err, jwt.ErrTokenMalformed) {
				return nil, fmt.Errorf("token is malformed")
			} else if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
				return nil, fmt.Errorf("token signature is invalid")
			} else if errors.Is(err, jwt.ErrTokenExpired) {
				return nil, fmt.Errorf("token is expired")
			} else if errors.Is(err, jwt.ErrTokenNotValidYet) {
				return nil, fmt.Errorf("token is not active yet")
			} else if errors.Is(err, jwt.ErrTokenUnverifiable) {
				return nil, fmt.Errorf("token could not be verified due to parsing issues")
			} else {
				return nil, fmt.Errorf("couldn't handle this token: %w", err)
			}
		}
		
		return nil, fmt.Errorf("invalid token")
	}

	// The claims are now in the claims variable
	return claims, nil
}
