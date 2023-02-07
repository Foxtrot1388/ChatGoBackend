package logging

import (
	logrus "github.com/sirupsen/logrus"
	"os"
)

var e *logrus.Entry

type Logger struct {
	*logrus.Entry
}

func init() {

	var log = &logrus.Logger{
		Out:       os.Stderr,
		Formatter: new(logrus.TextFormatter),
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.TraceLevel,
	}

	e = logrus.NewEntry(log)

}

func GetLogger() *Logger {
	return &Logger{e}
}
