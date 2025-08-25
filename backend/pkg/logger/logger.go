package logger

import (
	"log"
	"os"
)

type Logger struct {
	*log.Logger
}

func New() *Logger {
	return &Logger{
		Logger: log.New(os.Stdout, "PICKUP-QUEUE: ", log.LstdFlags|log.Lshortfile),
	}
}

func (l *Logger) Info(v ...interface{}) {
	l.Println("[INFO]", v)
}

func (l *Logger) Error(v ...interface{}) {
	l.Println("[ERROR]", v)
}

func (l *Logger) Warning(v ...interface{}) {
	l.Println("[WARNING]", v)
}

func (l *Logger) Debug(v ...interface{}) {
	l.Println("[DEBUG]", v)
}
