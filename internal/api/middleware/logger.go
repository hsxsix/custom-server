/**
 * @File: logger.go
 * @Author: hsien
 * @Description:
 * @Date: 9/18/21 3:02 PM
 */

package middleware

import (
	"bytes"
	"custom_server/pkg/log"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"time"
)

func Logger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//path := ctx.Request.URL.Path
		path := ctx.Request.URL.String()
		method := ctx.Request.Method

		log.Debug("Request ->",
			log.String("request id", requestid.Get(ctx)),
			log.String("client ip", ctx.ClientIP()),
			log.String("method", method),
			log.String("path", path),
			//log.Any("body", ctx.Request.Body),
		)
		start := time.Now()

		ctx.Next()

		latency := time.Now().Sub(start)
		statusCode := ctx.Writer.Status()
		msg := "Response ->"
		field := []log.Field{log.String("request id", requestid.Get(ctx)),
			//log.String("method", method),
			//log.String("path", path),
			log.Int("status code", ctx.Writer.Status()),
			log.String("time used", latency.String())}

		if statusCode >= http.StatusInternalServerError {
			log.Error(msg, field...)
		} else if statusCode >= http.StatusBadRequest {
			log.Warn(msg, field...)
		} else {
			log.Debug(msg, field...)
		}
	}
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func LoggerWithBody() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//path := ctx.Request.URL.Path
		path := ctx.Request.URL.String()
		method := ctx.Request.Method

		fields := []log.Field{
			log.String("request id", requestid.Get(ctx)),
			log.String("client ip", ctx.ClientIP()),
			log.String("method", method),
			log.String("path", path),
		}
		if len(GetRequestBody(ctx)) > 0 {
			fields = append(fields, log.Any("body", GetRequestBody(ctx)))
		}
		if len(ctx.Request.Form) > 0 {
			fields = append(fields, log.Any("form", ctx.Request.Form))
		}
		log.Debug("Request ->", fields...)
		start := time.Now()

		writer := &bodyLogWriter{body: bytes.NewBuffer([]byte{}), ResponseWriter: ctx.Writer}
		ctx.Writer = writer

		ctx.Next()

		latency := time.Now().Sub(start)
		statusCode := ctx.Writer.Status()
		msg := "Response ->"
		field := []log.Field{log.String("request id", requestid.Get(ctx)),
			//log.String("method", method),
			//log.String("path", path),
			log.Int("status code", ctx.Writer.Status()),
			log.String("time used", latency.String()),
			log.String("body", writer.body.String())}

		if statusCode >= http.StatusInternalServerError {
			log.Error(msg, field...)
		} else if statusCode >= http.StatusBadRequest {
			log.Warn(msg, field...)
		} else {
			log.Debug(msg, field...)
		}
	}
}

func GetRequestBody(c *gin.Context) []byte {
	// 获取请求 body
	var requestBody []byte
	if c.Request.Body != nil {
		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.Error(err)
		} else {
			requestBody = body
			// body 被 read 、 bind 之后会被置空，需要重置
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		}
	}
	return requestBody
}
