package function

import (
	"crypto/rsa"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

func createToken(email string, rsaKey string) (string, error) {
	claims := jwt.StandardClaims{}
	claims.Audience = email
	claims.Subject = "early-access"
	claims.ExpiresAt = time.Now().AddDate(0, 3, 0).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{})

	key, err := loadRSAPrivateKeyFromString(string(rsaKey))
	if err != nil {
		return "", fmt.Errorf("error creating private key from secret. %v", err)
	}

	return token.SignedString(key)
}

func loadRSAKeyFromSecret() ([]byte, error) {
	rsaKey, err := getSecret("faas-lambda-private-key")
	if err != nil {
		return nil, fmt.Errorf("error loading private key secret. %v", err)
	}

	return rsaKey, nil
}

func loadRSAPrivateKeyFromString(key string) (*rsa.PrivateKey, error) {
	return jwt.ParseRSAPrivateKeyFromPEM([]byte(key))
}