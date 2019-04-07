package credentials

import (
	"fmt"
	"log"

	"github.com/dgrijalva/jwt-go"
)

const (
	JWT_SECRET_KEY = "./credentials/secret_key"
	JWT_TOKEN_FILE = "./credentials/token.jwt"

	SSL_SERVER_CERT        = "./credentials/cert/server.crt"
	SSL_SERVER_PRIVATE_KEY = "./credentials/cert/server.key"
)

func ValidateToken(token string, secretKey string) (*jwt.Token, error) {
	jwtToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Printf("Unexpected signing method: %v", t.Header["alg"])
			return nil, fmt.Errorf("invalid token")
		}
		hmacSecret := []byte(secretKey)
		return hmacSecret, nil
	})
	if err == nil && jwtToken.Valid {
		return jwtToken, nil
	}
	return nil, err
}
