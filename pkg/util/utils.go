/**
 * @File: utils.go
 * @Author: hsien
 * @Description:
 * @Date: 9/24/21 9:40 AM
 */

package util

import (
	"context"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
)

func CtxWithRequestID(ctx *gin.Context) context.Context {
	return context.WithValue(context.Background(), "X-Request-ID", requestid.Get(ctx))
}
