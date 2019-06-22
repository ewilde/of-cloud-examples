package function

import (
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/openfaas/openfaas-cloud/sdk"
)

func createToken(email string, rsaKey string) (string, error) {
	claims := jwt.StandardClaims{}
	claims.Audience = email
	claims.Subject = "early-access"
	claims.ExpiresAt = time.Now().AddDate(0, 3, 0).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	key, err := loadRSAPrivateKeyFromString(string(rsaKey))
	if err != nil {
		return "", fmt.Errorf("error creating private key from secret. %v", err)
	}

	return token.SignedString(key)
}

func loadRSAKeyFromSecret() (string, error) {
	rsaKey, err := sdk.ReadSecret("faas-lambda-private.key")
	if err != nil {
		return "", fmt.Errorf("error loading private key secret. %v", err)
	}

	return rsaKey, nil
}

func loadRSAPrivateKeyFromString(key string) (*rsa.PrivateKey, error) {
	return jwt.ParseRSAPrivateKeyFromPEM([]byte(key))
}
