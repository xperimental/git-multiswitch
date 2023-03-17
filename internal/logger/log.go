package logger

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

type Logger interface {
	logrus.FieldLogger
	WriterLevel(level logrus.Level) *io.PipeWriter
}

func New() Logger {
	return &logrus.Logger{
		Out: os.Stderr,
		Formatter: &logrus.TextFormatter{
			DisableTimestamp: true,
		},
		Hooks:        logrus.LevelHooks{},
		Level:        logrus.DebugLevel,
		ExitFunc:     os.Exit,
		ReportCaller: false,
	}
}
