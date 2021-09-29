/**
 * @File: jwt.go
 * @Author: hsien
 * @Description:
 * @Date: 9/24/21 11:14 PM
 */

package util

import (
	jwt "github.com/dgrijalva/jwt-go"
)

//var jwtSecret = []byte("jwt_test")

type Claims struct {
	UserId string `json:"user_id"`
	jwt.StandardClaims
}

func GenerateToken(jwtSecret string, claims jwt.Claims) (string, error) {
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tokenClaims.SignedString([]byte(jwtSecret))
}

func ParseToken(jwtSecret, token string) (jwt.Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		if tokenClaims.Valid {
			return tokenClaims.Claims, nil
		}
	}

	return nil, err
}
