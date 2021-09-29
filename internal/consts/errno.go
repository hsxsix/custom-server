/**
 * @File: errno.go
 * @Author: hsien
 * @Description:
 * @Date: 9/18/21 4:23 PM
 */

package consts

import "custom_server/pkg/errno"

var (
	Success         = &errno.Errno{Code: 0, Msg: "ok"}
	InternalError   = &errno.Errno{Code: 10001, Msg: "internal server error"}
	ParamsBindError = &errno.Errno{Code: 10002, Msg: "parameter binding error"}

	PermissionDenied  = &errno.Err{Code: 10101, Msg: "permission denied"}
	TokenInvalid      = &errno.Err{Code: 10102, Msg: "token invalid"}
	TokenExpired      = &errno.Err{Code: 10103, Msg: "token expired"}
	PhoneRegistered   = &errno.Errno{Code: 10104, Msg: "phone number has been registered"}
	IncorrectPassword = &errno.Err{Code: 10104, Msg: "incorrect password"}
)

func BindErrWithError(err error) *errno.Err {
	return &errno.Err{Code: 10002, Msg: "parameter error", Err: err}
}

func TokenInvalidWithError(err error) *errno.Err {
	return &errno.Err{Code: 10102, Msg: "token invalid", Err: err}
}

func TokenExpiredWithError(err error) *errno.Err {
	return &errno.Err{Code: 10102, Msg: "token invalid", Err: err}
}

func IncorrectPasswordWithError(err error) *errno.Err {
	return &errno.Err{Code: 10104, Msg: "incorrect password", Err: err}
}
