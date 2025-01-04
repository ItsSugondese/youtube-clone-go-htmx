package utils

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func CreateToken(ttl time.Duration, payload interface{}, privateKey string) (string, error) {
	decodedPrivateKey, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		return "", fmt.Errorf("could not decode key: %w", err)
	}
	key, err := jwt.ParseRSAPrivateKeyFromPEM(decodedPrivateKey)

	if err != nil {
		return "", fmt.Errorf("create: parse keye: %w", err)
	}

	now := time.Now().UTC()

	claims := make(jwt.MapClaims)
	claims["sub"] = payload
	claims["exp"] = now.Add(ttl).Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)

	if err != nil {
		return "", fmt.Errorf("create: sign token: %w", err)
	}

	return token, nil
}

// ValidateToken validates the provided token using the given public key
func ValidateToken(token string, publicKeyBase64 string) bool {
	// Decode the base64-encoded public key
	decodedPublicKey, err := base64.StdEncoding.DecodeString(publicKeyBase64)
	if err != nil {
		panic(fmt.Sprintf("could not decode public key: %v", err))
	}

	// Parse the decoded public key
	key, err := jwt.ParseRSAPublicKeyFromPEM(decodedPublicKey)
	if err != nil {
		panic(fmt.Sprintf("validate: parse key: %w", err))
	}

	// Parse and validate the token
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			panic(fmt.Errorf("unexpected method: %s", t.Header["alg"]))
		}
		return key, nil
	})

	if err != nil {
		panic(fmt.Errorf("validate: %w", err))
	}

	// Check the token claims and validity
	_, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		panic("validate: invalid token")
	}

	return true
}
