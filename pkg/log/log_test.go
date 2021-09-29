/**
 * @File: log_test.go
 * @Author: hsien
 * @Description:
 * @Date: 9/16/21 5:13 PM
 */

package log

import "testing"

func TestLogger_Custom(t *testing.T) {
	WithOption(SetLevel("DEBUG"), ColorPrint(true), FileName("test.log"))
	Debug("this is a log with debug")
	Info("this is a log with info")
	Warn("this is a log with warn")
	Error("this is a log with error")
	//Panic("this is a log with panic")
	//Fatal("this is a log with fatal")
}

func TestLogger_Info(t *testing.T) {
	Info("this is a log with info level")
}

func TestLogger_Debug(t *testing.T) {
	Debug("this is a log with info debug")
}
