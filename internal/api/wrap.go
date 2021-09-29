/**
 * @File: wrap.go
 * @Author: hsien
 * @Description:
 * @Date: 9/18/21 4:30 PM
 */

package api

import (
	"custom_server/internal/consts"
	"custom_server/internal/model/response"
	"custom_server/pkg/config"
	"custom_server/pkg/errno"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler func(ctx *gin.Context) interface{}

func Wrap(handler Handler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var response = &response.Response{
			Code: consts.Success.Code,
			Msg:  consts.Success.Msg,
		}

		result := handler(ctx)
		switch t := result.(type) {
		case *errno.Errno:
			response.Code = t.Code
			response.Msg = t.Msg
		case *errno.Err:
			response.Code = t.Code
			response.Msg = t.Msg
			if config.Debug {
				response.Msg = t.Err.Error()
			}
		case error:
			response.Code = consts.InternalError.Code
			response.Msg = consts.InternalError.Msg
			if config.Debug {
				response.Msg = t.Error()
			}
		default:
			response.Data = result
		}
		ctx.JSON(http.StatusOK, response)
	}
}
