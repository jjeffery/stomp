package stomp

import "github.com/go-stomp/stomp/v3/internal/log"

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

// SetLogger sets the logger used. You need to call this method
// only if you want to customise the logging from stomp package.
func SetLogger(l Logger) {
	log.SetLogger(l)
}
