package utils

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSuperSecretKey = []byte("supersecretkey")

func VerifyPassword(password, hashedPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil, err
}

func GenerateTokens(username, role string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"role":     role,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSuperSecretKey)

}

func VerifyTokens(tokenString string) (*jwt.MapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		// right now just return the secret key
		return jwtSuperSecretKey, nil
	})

	// check the validity of the token or check for errors
	if err != nil || !token.Valid {
		return nil, errors.New("unauthorized: invalid token")
	}

	// we need to extract the claims
	claims, ok := token.Claims.(*jwt.MapClaims)
	if !ok {
		return nil, errors.New("unauthorized: invalid claims")
	}

	return claims, nil

}
