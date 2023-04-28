package jwttoken

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

func CreateToken(useremail string) (string, error) {

	expirationTime := time.Now().Add(time.Minute * 5)

	claims := &Claims{
		UserEmail: useremail,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			// ExpiresAt: time.Now().Add(time.Duration(tokenLifeTime)).Unix(),
		},
	}

	// create token using Claims and jwtKey

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil

}
