package log

import (
	"os"
)

type Logger interface {
	Debugf(format string, value ...interface{})
	Infof(format string, value ...interface{})
	Warningf(format string, value ...interface{})
	Errorf(format string, value ...interface{})

	Debug(message string)
	Info(message string)
	Warning(message string)
	Error(message string)
}

var logger Logger

func SetLogger(l Logger) {
	if l != nil {
		logger = l
	}
}

func Debugf(format string, value ...interface{}) {
	logger.Debugf(format, value...)
}

func Debug(message string) {
	Debugf("%s", message)
}

func Infof(format string, value ...interface{}) {
	logger.Infof(format, value...)
}

func Info(message string) {
	Infof("%s", message)
}

func Warningf(format string, value ...interface{}) {
	logger.Warningf(format, value...)
}

func Warning(message string) {
	Warningf("%s", message)
}

func Errorf(format string, value ...interface{}) {
	logger.Errorf(format, value...)
}

func Error(message string) {
	Errorf("%s", message)
}

func Fatalf(format string, value ...interface{}) {
	Errorf(format, value...)
	os.Exit(1)
}

func Fatal(message string) {
	Fatalf("%s", message)
}

func init() {
	logger = stdLogger{}
}
