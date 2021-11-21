package util

import (
	"blog/models"
	"fmt"
	"os"

	"github.com/golang-jwt/jwt"
)

func CheckToken(userEmail string, tokenString string) bool {
	if tokenString == "" {
		return false
	}
	fmt.Println("获得的token: ", tokenString)
	var hmacSampleSecret []byte = []byte(os.Getenv("TOKEN_SALT"))

	token, err := jwt.ParseWithClaims(tokenString, &models.MyClaims{}, func(t *jwt.Token) (interface{}, error) {
		return hmacSampleSecret, nil
	})
	if err != nil {
		fmt.Println("token 解析错误:", err)
		return false
	}
	return token.Claims.(*models.MyClaims).Email == userEmail
}
