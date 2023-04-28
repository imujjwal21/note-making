package jwttoken

import "github.com/dgrijalva/jwt-go"

var jwtKey = []byte("secret_key") // it is jwt secret key that use in jwt tokento sign our jwt token / used to create the signature

type Claims struct {
	UserEmail string `json:"Useremail"`
	jwt.StandardClaims
}
