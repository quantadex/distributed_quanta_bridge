package logger

import "github.com/op/go-logging"

type GoLogger struct {
	log *logging.Logger
}

var Formatter = logging.MustStringFormatter("%{level:.1s} [%{module}] %{message}")

func NewGoLogger(module string ) *GoLogger {
	logging.SetFormatter(Formatter)
	var log = logging.MustGetLogger(module)

	return &GoLogger{log}
}

func (l *GoLogger) Error(msg string) {
	l.log.Error(msg)
}

func (l *GoLogger) Info(msg string) {
	l.log.Info(msg)
}

func (l *GoLogger) Infof(format string, args ...interface{}) {
	l.log.Infof(format, args...)
}

func (l *GoLogger) Debug(msg string) {
	l.log.Debug(msg)
}



