/**
 * @File: errno.go
 * @Author: hsien
 * @Description:
 * @Date: 9/18/21 4:19 PM
 */

package errno

import "fmt"

var (
	Success       = &Errno{Code: 0, Msg: "ok"}
	InternalError = &Errno{Code: 10001, Msg: "internal server error"}
)

type Errno struct {
	Code int
	Msg  string
}

func (err Errno) Error() string {
	return err.Msg
}

type Err struct {
	Code int
	Msg  string
	Err  error
}

func (err Err) Error() string {
	if err.Err != nil {
		return fmt.Sprintf("err - code: %d, msg: %s, error:%s", err.Code, err.Msg, err.Err)
	}
	return fmt.Sprintf("err - code: %d, error:%s", err.Code, err.Msg)
}
