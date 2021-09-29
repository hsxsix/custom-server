/**
 * @File: user.go
 * @Author: hsien
 * @Description: response model
 * @Date: 9/22/21 5:49 PM
 */

package response

// Response params
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// Page params
type Page struct {
	PageNo   int `json:"page_no" form:"page_no"`
	PageSize int `json:"page_size" form:"page_size"`
}

type UserInfoRes struct {
	UserID      string `json:"user_id"`
	UserName    string `json:"user_name"`
	PhoneNumber string `json:"phone_number"`
	UserAvatar  string `json:"user_avatar"`
	UserSex     string `json:"user_sex"`
	UserDesc    string `json:"user_desc"`
	CreatedAt   string `json:"created_at"`
}

type UserListRes struct {
	Total int `json:"total"`
	Page
	List []*UserInfoRes `json:"list"`
}

type UserTokenRes struct {
	Token string `json:"token"`
}
