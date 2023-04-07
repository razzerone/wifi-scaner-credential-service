package autorisation

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// SignedToken Формирование JWT.
func SignedToken(key *Key) string {

	claims := jwt.RegisteredClaims{
		Issuer:    key.AccountID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		Audience:  []string{"https://iam.api.cloud.yandex.net/iam/v1/tokens"},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodPS256, claims)
	token.Header["kid"] = key.ID

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(key.PrivateKey))
	if err != nil {
		log.Panic("unable to parse private key")
	}

	signed, err := token.SignedString(privateKey)
	if err != nil {
		panic(err)
	}
	return signed
}
