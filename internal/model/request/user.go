/**
 * @File: user.go
 * @Author: hsien
 * @Description: request model
 * @Date: 9/23/21 3:41 PM
 */

package request

import (
	"custom_server/internal/model/response"
)

type UserQuery struct {
	response.Page
	Name        string `form:"username"`
	Sex         int    `form:"sex"`
	PhoneNumber string `form:"phone_number"`
}

func DefaultUserQuery() *UserQuery {
	query := new(UserQuery)
	query.PageNo = 1
	query.PageSize = 10
	return query
}

type UserRegister struct {
	UserName    string `json:"username"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type UserLogin struct {
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}
