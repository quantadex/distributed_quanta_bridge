package logger

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
}

func NewLogger() (Logger, error) {
    return &GoLogger{}, nil
}
