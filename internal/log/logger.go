package log

import (
	"os"

	"github.com/go-stomp/stomp/v3"
)

var logger stomp.Logger

func SetLogger(l stomp.Logger) {
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
