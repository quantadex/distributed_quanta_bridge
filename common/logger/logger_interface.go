package logger

import "github.com/op/go-logging"

/**
 * Logger
 *
 * Used by all modules. Adds datetime and file/line to log lines.
 * Additionally may push logs to another location.
 */
type Logger interface {
    Error(msg string)
    Info(msg string)
    Debug(msg string)
    Errorf(format string, args ...interface{})
    Infof(format string, args ...interface{})
    Debugf(format string, args ...interface{})

    SetLogLevel(level logging.Level)
}

func NewLogger(module string) (Logger, error) {
    return NewGoLogger(module), nil
}
