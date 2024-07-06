package helpers

import (
	"crypto/x509"
	"encoding/pem"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(clientId, privatePem string) (string, error) {
	block, _ := pem.Decode([]byte(privatePem))
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return "", nil
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}

	now := time.Now().Unix()
	payload := jwt.MapClaims{
		"iat": now - 60,
		"exp": now + (10 * 60),
		"iss": clientId,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, payload)

	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
