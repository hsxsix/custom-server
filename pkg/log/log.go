/**
 * @File: log.go
 * @Author: hsien
 * @Description: learn(copy) from https://tonybai.com/2021/07/14/uber-zap-advanced-usage/
 * @Date: 9/16/21 3:18 PM
 */

package log

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"strings"
)

type Level = zapcore.Level

const (
	DebugLevel  Level = zap.DebugLevel
	InfoLevel   Level = zap.InfoLevel
	WarnLevel   Level = zap.WarnLevel
	ErrorLevel  Level = zap.ErrorLevel
	DPanicLevel Level = zap.DPanicLevel
	PanicLevel  Level = zap.PanicLevel
	FatalLevel  Level = zap.FatalLevel
)

type Field = zapcore.Field

type ZapOption = zap.Option

var (
	WithCaller = zap.WithCaller
	// AddCallerSkip AddStacktrace = zap.AddStacktrace
	AddCallerSkip = zap.AddCallerSkip
)

var tops = []TeeOption{
	//{
	//   FileName: "log.log",
	//   Lel: DebugLevel,
	//},
	{
		W:          os.Stderr,
		Lel:        DebugLevel,
		ColorPrint: true,
	},
}

var std = newLogger(tops, WithCaller(true), AddCallerSkip(1))

var (
	Skip        = zap.Skip
	Binary      = zap.Binary
	Bool        = zap.Bool
	Boolp       = zap.Boolp
	ByteString  = zap.ByteString
	Complex128  = zap.Complex128
	Complex128p = zap.Complex128p
	Complex64   = zap.Complex64
	Complex64p  = zap.Complex64p
	Float64     = zap.Float64
	Float64p    = zap.Float64p
	Float32     = zap.Float32
	Float32p    = zap.Float32p
	Int         = zap.Int
	Intp        = zap.Intp
	Int64       = zap.Int64
	Int64p      = zap.Int64p
	Int32       = zap.Int32
	Int32p      = zap.Int32p
	Int16       = zap.Int16
	Int16p      = zap.Int16p
	Int8        = zap.Int8
	Int8p       = zap.Int8p
	String      = zap.String
	Stringp     = zap.Stringp
	Uint        = zap.Uint
	Uintp       = zap.Uintptrp
	Uint64      = zap.Uint64
	Uint64p     = zap.Uint64p
	Uint32      = zap.Uint32
	Uint32p     = zap.Uint32p
	Uint16      = zap.Uint16
	Uint16p     = zap.Uint16p
	Uint8       = zap.Uint8
	Uint8p      = zap.Uint8p
	Uintptr     = zap.Uintptr
	Uintptrp    = zap.Uintptrp
	Reflect     = zap.Reflect
	NameError   = zap.NamedError
	Namespace   = zap.Namespace
	Stringer    = zap.Stringer
	Time        = zap.Time
	Timep       = zap.Timep
	Stack       = zap.Stack
	StackSkip   = zap.StackSkip
	Duration    = zap.Duration
	Durationp   = zap.Durationp
	Object      = zap.Object
	Inline      = zap.Inline
	Any         = zap.Any

	Debug  = std.Debug
	Debugf = std.Debugf
	Info   = std.Info
	Infof  = std.Infof
	Warn   = std.Warn
	Warnf  = std.Warnf
	Error  = std.Error
	Errorf = std.Errorf
	DPanic = std.DPanic
	Panic  = std.Panic
	Panicf = std.Panicf
	Fatal  = std.Fatal
	Fatalf = std.Fatalf
)

type Logger struct {
	logger *zap.Logger
	sugar  *zap.SugaredLogger
	level  Level
}

type LevelEnableFunc func(lvl Level) bool

type TeeOption struct {
	FileName   string
	W          io.Writer
	Lef        LevelEnableFunc
	Lel        Level
	ColorPrint bool
}

func newLogger(tops []TeeOption, opts ...ZapOption) *Logger {
	var cores []zapcore.Core

	cfg := zap.NewProductionConfig().EncoderConfig
	cfg.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")

	for _, top := range tops {
		top := top

		lv := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= top.Lel
		})
		if top.Lef != nil {
			lv = func(lvl zapcore.Level) bool {
				return top.Lef(lvl)
			}
		}

		if top.FileName != "" {
			w := zapcore.AddSync(&lumberjack.Logger{
				Filename: top.FileName,
				MaxSize:  100,
				MaxAge:   15,
				//MaxBackups: 3,
				Compress: true,
			})
			top.W = w
		}

		if top.W == os.Stderr {
			stdCfg := cfg
			if top.ColorPrint {
				stdCfg.EncodeLevel = zapcore.LowercaseColorLevelEncoder
			}
			cores = append(cores, zapcore.NewCore(
				zapcore.NewConsoleEncoder(stdCfg),
				zapcore.AddSync(top.W),
				lv,
			))
		} else {
			cores = append(cores, zapcore.NewCore(
				zapcore.NewJSONEncoder(cfg),
				zapcore.AddSync(top.W),
				lv,
			))
		}
	}

	logger := &Logger{
		logger: zap.New(zapcore.NewTee(cores...), opts...),
	}
	logger.sugar = logger.logger.Sugar()
	return logger
}

func (l *Logger) Sync() error {
	var err error
	err = l.logger.Sync()
	err = l.sugar.Sync()
	return err
}

func (l *Logger) Debug(msg string, fields ...Field) {
	l.logger.Debug(msg, fields...)
}

func (l *Logger) Debugf(template string, args ...interface{}) {
	l.sugar.Debugf(template, args...)
}

func (l *Logger) Info(msg string, fields ...Field) {
	l.logger.Info(msg, fields...)
}

func (l *Logger) Infof(template string, args ...interface{}) {
	l.sugar.Infof(template, args...)
}

func (l *Logger) Warn(msg string, fields ...Field) {
	l.logger.Warn(msg, fields...)
}

func (l *Logger) Warnf(template string, args ...interface{}) {
	l.sugar.Warnf(template, args...)
}

func (l *Logger) Error(msg string, fields ...Field) {
	l.logger.Error(msg, fields...)
}

func (l *Logger) Errorf(template string, args ...interface{}) {
	l.sugar.Errorf(template, args...)
}

func (l *Logger) DPanic(msg string, fields ...Field) {
	l.logger.DPanic(msg, fields...)
}

func (l *Logger) Panic(msg string, fields ...Field) {
	l.logger.Panic(msg, fields...)
}

func (l *Logger) Panicf(template string, args ...interface{}) {
	l.sugar.Panicf(template, args...)
}

func (l *Logger) Fatal(msg string, fields ...Field) {
	l.logger.Fatal(msg, fields...)
}

func (l *Logger) Fatalf(template string, args ...interface{}) {
	l.sugar.Fatalf(template, args...)
}

func reset(l *Logger) {
	std = l
	Debug = std.Debug
	Debugf = std.Debugf
	Info = std.Info
	Infof = std.Infof
	Warn = std.Warn
	Warnf = std.Warnf
	Error = std.Error
	Errorf = std.Errorf
	DPanic = std.DPanic
	Panic = std.Panic
	Panicf = std.Panicf
	Fatal = std.Fatal
	Fatalf = std.Fatalf
	DebugWithGinCtx = std.DebugWithGinCtx
	InfoWithGinCtx = std.InfoWithGinCtx
	WarnWithGinCtx = std.WarnWithGinCtx
	ErrorWithGinCtx = std.ErrorWithGinCtx
	PanicWithGinCtx = std.PanicWithGinCtx
	FatalWithGinCtx = std.FatalWithGinCtx
	DebugWithCtx = std.DebugWithCtx
	InfoWithCtx = std.InfoWithCtx
	WarnWithCtx = std.WarnWithCtx
	ErrorWithCtx = std.ErrorWithCtx
	PanicWithCtx = std.PanicWithCtx
	FatalWithCtx = std.FatalWithCtx
}

func Sync() error {
	if std != nil {
		return std.Sync()
	}
	return nil
}

type config struct {
	fileName   string
	level      Level
	colorPrint bool
}

type option func(*config)

// SetLevel set log level
func SetLevel(lvl string) option {
	return func(c *config) {
		c.level = ParseLevel(lvl)
	}
}

func ParseLevel(lvl string) Level {
	var level Level
	switch strings.ToUpper(lvl) {
	case "DEBUG":
		level = DebugLevel
	case "INFO":
		level = InfoLevel
	case "WARN":
		level = WarnLevel
	case "ERROR":
		level = ErrorLevel
	case "PANIC":
		level = PanicLevel
	case "FATAL":
		level = FatalLevel
	default:
		level = DebugLevel // default Debug
	}
	return level
}

// FileName log file path(name), if is empty, Will not be recorded to file
func FileName(name string) option {
	return func(c *config) {
		c.fileName = name
	}
}

// ColorPrint set color print, only for console
func ColorPrint(b bool) option {
	return func(c *config) {
		c.colorPrint = b
	}
}

func WithOption(opts ...option) {
	if len(opts) == 0 {
		return
	}
	cfg := config{
		level: DebugLevel,
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	var tops []TeeOption
	tops = append(tops, TeeOption{
		W:          os.Stderr,
		Lel:        cfg.level,
		ColorPrint: cfg.colorPrint,
	})
	if cfg.fileName != "" {
		tops = append(tops, TeeOption{
			FileName: cfg.fileName,
			Lel:      cfg.level,
		})
	}
	reset(newLogger(tops, WithCaller(true), AddCallerSkip(1)))
}
