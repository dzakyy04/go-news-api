package utils

import (
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
