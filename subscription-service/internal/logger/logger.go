package logger

import (
	"log"
	"os"
)

type Logger struct {
	infoLog  *log.Logger
	errorLog *log.Logger
}

// NewLogger creates a new logger instance
func NewLogger(level string) *Logger {
	// level parameter accepted for future use, currently not filtering
	return &Logger{
		infoLog:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		errorLog: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func (l *Logger) Info(v ...interface{}) {
	l.infoLog.Println(v...)
}

func (l *Logger) Error(v ...interface{}) {
	l.errorLog.Println(v...)
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.infoLog.Printf(format, v...)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.errorLog.Printf(format, v...)
}
