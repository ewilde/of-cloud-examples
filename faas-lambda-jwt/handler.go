package function

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Handle a serverless request
func Handle(req []byte) string {
	email := string(req)
	claims := jwt.StandardClaims{}
	claims.Audience = email
	claims.Subject = "early-access"
	claims.ExpiresAt = time.Now().AddDate(0, 3, 0).Unix()
	jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{})

	key, err := loadRSAKeyFromSecret()
	if err != nil {
		return err.Error()
	}

	token, err := createToken(email, string(key))
	if err != nil {
		return err.Error()
	}

	return token
}
