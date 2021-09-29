/**
 * @File: GormLogger.go
 * @Author: hsien
 * @Description:
 * @Date: 9/24/21 2:38 PM
 */

package log

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm/logger"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type GormLogger struct {
	GormLogConfig
}

type GormLogConfig struct {
	SlowThreshold             time.Duration
	IgnoreRecordNotFoundError bool
	LogLevel                  Level
}

func (g *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	gl := *g
	switch level {
	case logger.Silent:
		gl.LogLevel = PanicLevel
	case logger.Error:
		gl.LogLevel = ErrorLevel
	case logger.Warn:
		gl.LogLevel = WarnLevel
	case logger.Info:
		gl.LogLevel = DebugLevel
	default:
		gl.LogLevel = WarnLevel
	}
	return &gl
}

func (g GormLogger) Info(ctx context.Context, format string, args ...interface{}) {
	if g.LogLevel > InfoLevel {
		return
	}
	requestId, shortPath := GetRequestId(ctx), g.GetShortPath()
	var fields []Field
	if requestId != "" {
		fields = append(fields, String("request id", requestId))
	}
	fields = append(fields, String("line", shortPath), String("info", fmt.Sprintf(format, args...)))
	Info("Gorm ->", fields...)
}

func (g GormLogger) Warn(ctx context.Context, format string, args ...interface{}) {
	if g.LogLevel > WarnLevel {
		return
	}
	requestId, shortPath := GetRequestId(ctx), g.GetShortPath()
	var fields []Field
	if requestId != "" {
		fields = append(fields, String("request id", requestId))
	}
	fields = append(fields, String("line", shortPath), String("info", fmt.Sprintf(format, args...)))
	Warn("Gorm ->", fields...)
}

func (g GormLogger) Error(ctx context.Context, format string, args ...interface{}) {
	if g.LogLevel > ErrorLevel {
		return
	}
	requestId, shortPath := GetRequestId(ctx), g.GetShortPath()
	var fields []Field
	if requestId != "" {
		fields = append(fields, String("request id", requestId))
	}
	fields = append(fields, String("line", shortPath), String("info", fmt.Sprintf(format, args...)))
	Error("Gorm ->", fields...)
}

func (g GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	requestId, elapsed := GetRequestId(ctx), time.Since(begin)
	shortPath := g.GetShortPath()
	switch {
	// print error log
	case err != nil && g.LogLevel <= ErrorLevel && (!errors.Is(err, logger.ErrRecordNotFound) || !g.IgnoreRecordNotFoundError):
		sql, rows := fc()
		var fields []Field
		if requestId != "" {
			fields = append(fields, String("request id", requestId))
		}
		fields = append(fields,
			String("line", shortPath),
			NameError("error", err),
			String("time", elapsed.String()))
		if rows == -1 {
			fields = append(fields, String("rows", "-"))
		} else {
			fields = append(fields, Int64("rows", rows))
		}
		fields = append(fields, String("sql", sql))
		Error("Gorm ->", fields...)
	// print warn log
	case elapsed > g.SlowThreshold && g.SlowThreshold != 0 && g.LogLevel <= WarnLevel:
		sql, rows := fc()
		var fields []Field
		if requestId != "" {
			fields = append(fields, String("request id", requestId))
		}
		fields = append(fields,
			String("line", shortPath),
			String("slow sql >=", g.SlowThreshold.String()),
			String("time", elapsed.String()))
		if rows == -1 {
			fields = append(fields, String("rows", "-"))
		} else {
			fields = append(fields, Int64("rows", rows))
		}
		fields = append(fields, String("sql", sql))
		Warn("Gorm ->", fields...)
	// print debug log
	case g.LogLevel < InfoLevel:
		sql, rows := fc()
		var fields []Field
		if requestId != "" {
			fields = append(fields, String("request id", requestId))
		}
		fields = append(fields,
			String("line", shortPath),
			String("time", elapsed.String()))
		if rows == -1 {
			fields = append(fields, String("rows", "-"))
		} else {
			fields = append(fields, Int64("rows", rows))
		}

		fields = append(fields, String("sql", sql))
		Debug("Gorm ->", fields...)
	}
}

func (g GormLogger) GetShortPath() string {
	var path string
	// caller skip 4
	_, file, line, ok := runtime.Caller(4)
	if !ok {
		return ""
	}
	path = file + ":" + strconv.FormatInt(int64(line), 10)
	pathSplit := strings.Split(path, "/")
	l := len(pathSplit)
	if l > 2 {
		return strings.Join(pathSplit[l-2:l], "/")
	}
	return path
}

func NewGormLogger(logConfig GormLogConfig) *GormLogger {
	return &GormLogger{
		GormLogConfig: logConfig,
	}
}
