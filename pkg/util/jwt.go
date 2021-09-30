/**
 * @File: jwt.go
 * @Author: hsien
 * @Description:
 * @Date: 9/24/21 11:14 PM
 */

package util

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
)

//var jwtSecret = []byte("jwt_test")

func GenerateToken(jwtSecret string, claims jwt.Claims) (string, error) {
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tokenClaims.SignedString([]byte(jwtSecret))
}

func ParseToken(token, jwtSecret string, claims jwt.Claims) (jwt.Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !tokenClaims.Valid {
		return nil, errors.New("token expired")
	}
	return tokenClaims.Claims, nil
}
