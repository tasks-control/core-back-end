package utils

import (
	"fmt"
	"runtime"

	"github.com/sirupsen/logrus"
)

type LoggerInterface interface {
	logrus.FieldLogger
}

type logger struct {
	*logrus.Entry
}

var log *logger

func Logger() LoggerInterface {
	if log == nil {
		log = NewLogger().(*logger)
	}
	return log
}

func NewLogger() LoggerInterface {
	baseLogger := logrus.New()
	baseLogger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat:   "2006-01-02 15:04:05",
		DisableTimestamp:  false,
		DisableHTMLEscape: false,
		DataKey:           "data",
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  ".t",
			logrus.FieldKeyMsg:   "@msg",
			logrus.FieldKeyLevel: "@level",
		},
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			return frame.Function, fmt.Sprintf("%s:%d", frame.File, frame.Line)
		},
		PrettyPrint: false,
	})

	return &logger{logrus.NewEntry(baseLogger)}
}
