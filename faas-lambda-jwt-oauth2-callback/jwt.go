package function

import (
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/openfaas/openfaas-cloud/sdk"
)

const issuer = `https://ewilde.o6s.io/faas-lambda-jwt`
const audience = `faas-lambda`

func createToken(email string) (string, error) {
	keyPEM, err := readRSAKeyFromSecret()
	if err != nil {
		return "", fmt.Errorf("error reading secret. %v", err)
	}

	key, err := loadRSAPrivateKeyFromPEM(string(keyPEM))
	if err != nil {
		return "", fmt.Errorf("error creating private key from secret. %v", err)
	}

	claims := jwt.StandardClaims{}
	claims.Issuer = issuer
	claims.Audience = audience
	claims.Subject = email
	claims.ExpiresAt = time.Now().AddDate(0, 3, 0).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	return token.SignedString(key)
}

func readRSAKeyFromSecret() (string, error) {
	rsaKey, err := sdk.ReadSecret("faas-lambda-private.key")
	if err != nil {
		return "", fmt.Errorf("error loading private key secret. %v", err)
	}

	return rsaKey, nil
}

func loadRSAPrivateKeyFromPEM(key string) (*rsa.PrivateKey, error) {
	return jwt.ParseRSAPrivateKeyFromPEM([]byte(key))
}
