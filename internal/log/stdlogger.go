package log

import stdlog "log"

var (
	debugPrefix = "DEBUG: "
	infoPrefix  = "INFO: "
	warnPrefix  = "WARN: "
	errorPrefix = "ERROR: "
)

type stdLogger struct{}

func (s stdLogger) Debugf(format string, value ...interface{}) {
	stdlog.Printf(debugPrefix+format+"\n", value...)
}

func (s stdLogger) Debug(message string) {
	s.Debugf("%s", message)
}

func (s stdLogger) Infof(format string, value ...interface{}) {
	stdlog.Printf(infoPrefix+format+"\n", value...)
}

func (s stdLogger) Info(message string) {
	s.Infof("%s", message)
}

func (s stdLogger) Warningf(format string, value ...interface{}) {
	stdlog.Printf(warnPrefix+format+"\n", value...)
}

func (s stdLogger) Warning(message string) {
	s.Warningf("%s", message)
}

func (s stdLogger) Errorf(format string, value ...interface{}) {
	stdlog.Printf(errorPrefix+format+"\n", value...)
}

func (s stdLogger) Error(message string) {
	s.Errorf("%s", message)
}
