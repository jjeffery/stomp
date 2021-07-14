package log

import stdlog "log"

var (
	debugPrefix = "DEBUG: "
	infoPrefix  = "INFO: "
	warnPrefix  = "WARN: "
	errorPrefix = "ERROR: "
)

type StdLogger struct{}

func (s StdLogger) Debugf(format string, value ...interface{}) {
	stdlog.Printf(debugPrefix+format+"\n", value...)
}

func (s StdLogger) Debug(message string) {
	s.Debugf("%s", message)
}

func (s StdLogger) Infof(format string, value ...interface{}) {
	stdlog.Printf(infoPrefix+format+"\n", value...)
}

func (s StdLogger) Info(message string) {
	s.Infof("%s", message)
}

func (s StdLogger) Warningf(format string, value ...interface{}) {
	stdlog.Printf(warnPrefix+format+"\n", value...)
}

func (s StdLogger) Warning(message string) {
	s.Warningf("%s", message)
}

func (s StdLogger) Errorf(format string, value ...interface{}) {
	stdlog.Printf(errorPrefix+format+"\n", value...)
}

func (s StdLogger) Error(message string) {
	s.Errorf("%s", message)
}
