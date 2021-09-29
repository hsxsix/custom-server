/**
 * @File: withctx.go
 * @Author: hsien
 * @Description:
 * @Date: 9/24/21 9:49 AM
 */

package log

import (
	"context"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
)

var (
	DebugWithGinCtx = std.DebugWithGinCtx
	InfoWithGinCtx  = std.InfoWithGinCtx
	WarnWithGinCtx  = std.WarnWithGinCtx
	ErrorWithGinCtx = std.ErrorWithGinCtx
	PanicWithGinCtx = std.PanicWithGinCtx
	FatalWithGinCtx = std.FatalWithGinCtx
	DebugWithCtx    = std.DebugWithCtx
	InfoWithCtx     = std.InfoWithCtx
	WarnWithCtx     = std.WarnWithCtx
	ErrorWithCtx    = std.ErrorWithCtx
	PanicWithCtx    = std.PanicWithCtx
	FatalWithCtx    = std.FatalWithCtx
)

func (l *Logger) DebugWithGinCtx(ctx *gin.Context, msg string, fields ...Field) {
	l.logger.With(String("request id", requestid.Get(ctx))).Debug(msg, fields...)
}

func (l *Logger) InfoWithGinCtx(ctx *gin.Context, msg string, fields ...Field) {
	l.logger.With(String("request id", requestid.Get(ctx))).Info(msg, fields...)
}

func (l *Logger) WarnWithGinCtx(ctx *gin.Context, msg string, fields ...Field) {
	l.logger.With(String("request id", requestid.Get(ctx))).Warn(msg, fields...)
}

func (l *Logger) ErrorWithGinCtx(ctx *gin.Context, msg string, fields ...Field) {
	l.logger.With(String("request id", requestid.Get(ctx))).Error(msg, fields...)
}

func (l *Logger) PanicWithGinCtx(ctx *gin.Context, msg string, fields ...Field) {
	l.logger.With(String("request id", requestid.Get(ctx))).Panic(msg, fields...)
}

func (l *Logger) FatalWithGinCtx(ctx *gin.Context, msg string, fields ...Field) {
	l.logger.With(String("request id", requestid.Get(ctx))).Fatal(msg, fields...)
}

func (l *Logger) DebugWithCtx(ctx context.Context, msg string, fields ...Field) {
	l.logger.With(String("request id", GetRequestId(ctx))).Debug(msg, fields...)
}

func (l *Logger) InfoWithCtx(ctx context.Context, msg string, fields ...Field) {
	l.logger.With(String("request id", GetRequestId(ctx))).Info(msg, fields...)
}

func (l *Logger) WarnWithCtx(ctx context.Context, msg string, fields ...Field) {
	l.logger.With(String("request id", GetRequestId(ctx))).Warn(msg, fields...)
}

func (l *Logger) ErrorWithCtx(ctx context.Context, msg string, fields ...Field) {
	l.logger.With(String("request id", GetRequestId(ctx))).Error(msg, fields...)
}

func (l *Logger) PanicWithCtx(ctx context.Context, msg string, fields ...Field) {
	l.logger.With(String("request id", GetRequestId(ctx))).Panic(msg, fields...)
}

func (l *Logger) FatalWithCtx(ctx context.Context, msg string, fields ...Field) {
	l.logger.With(String("request id", GetRequestId(ctx))).Fatal(msg, fields...)
}

func GetRequestId(ctx context.Context) string {
	requestId, _ := ctx.Value("X-Request-ID").(string)
	return requestId
}
