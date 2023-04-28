package jwttoken

import (
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

func CheckToken(r *http.Request) string {

	cookie, err := r.Cookie("token")

	if err != nil {
		if err == http.ErrNoCookie {
			log.Printf("StatusUnauthorized Error : %v \n", err)
			return ""
		}
		log.Printf("StatusBadRequest Error : %v \n", err)
		return ""
	}

	tknStr := cookie.Value

	// Initialize a new instance of `Claims`
	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			log.Printf("StatusUnauthorized ParseWithClaim Error : %v \n", err)
			return ""
		}
		log.Printf("StatusBadRequest ParseWithClaim Error : %v \n", err)
		return ""
	}
	if !tkn.Valid {
		log.Printf("StatusUnauthorized Invalid Token Error : %v \n", err)
		return ""
	}

	return claims.UserEmail
}
