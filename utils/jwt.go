package utils

import (
	"errors"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte(os.Getenv("JWT_KEY"))

func GenerateToken(claims *jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwt, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return jwt, nil
}

func ParseToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtKey, nil
	})
}
