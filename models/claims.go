package models

import "github.com/golang-jwt/jwt"

type MyClaims struct {
	Email string `json:"username"`
	jwt.StandardClaims
}
