package logger

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

type Logger interface {
	logrus.FieldLogger
	SetLevel(level logrus.Level)
	WriterLevel(level logrus.Level) *io.PipeWriter
}

func New() Logger {
	return &logrus.Logger{
		Out: os.Stderr,
		Formatter: &logrus.TextFormatter{
			DisableTimestamp: true,
		},
		Hooks:        logrus.LevelHooks{},
		Level:        logrus.InfoLevel,
		ExitFunc:     os.Exit,
		ReportCaller: false,
	}
}
