package logger

import "github.com/op/go-logging"

type GoLogger struct {

}

var log = logging.MustGetLogger("bridge")

func (l *GoLogger) Error(msg string) {
	log.Error(msg)
}

func (l *GoLogger) Info(msg string) {
	log.Info(msg)
}

func (l *GoLogger) Debug(msg string) {
	log.Debug(msg)
}



