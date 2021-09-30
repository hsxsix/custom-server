/**
 * @File: jwt.go
 * @Author: hsien
 * @Description:
 * @Date: 9/25/21 12:53 AM
 */

package model

import "github.com/dgrijalva/jwt-go"

type UserClaims struct {
	UserId string `json:"user_id"`
	jwt.StandardClaims
}
