package auth

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

var _ Auth = (*JWT)(nil)

type JWT struct {
	Secret string
}

func NewJWT(secret string) JWT {
	return JWT{secret}
}

func (j JWT) GenerateToken() (string, error) {
	claims := jwt.MapClaims{
		"admin": true,
		"exp":   time.Now().Add(time.Hour).Unix(),
		"name":  "Ruslan",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(j.Secret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (j JWT) VerifyToken(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(j.Secret), nil
	})
	if err != nil {
		return false, err
	}

	return token.Valid, nil
}
