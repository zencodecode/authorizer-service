package rs256

import (
	"crypto/rsa"
	"errors"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateAccessToken[T jwt.Claims](privateKey *rsa.PrivateKey, claims T, kid string) (string, error) {

	token := jwt.NewWithClaims(
		jwt.SigningMethodRS256,
		claims,
	)
	token.Header["kid"] = kid
	return token.SignedString(privateKey)
}

func ValidateAccessToken[T jwt.Claims](tokenString string, publicKey *rsa.PublicKey, claims T) (T, error) {
	parsedToken, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return publicKey, nil
	})

	if err != nil {
		return claims, err
	}

	if !parsedToken.Valid {
		return claims, errors.New("token invalid")
	}
	return claims, nil
}
