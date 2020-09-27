package mocklogger

import (
	"github.com/getumen/replicator/pkg/logger"
	"github.com/sirupsen/logrus"
)

type stderrLogger struct {
	internal *logrus.Entry
}

// NewLogger returns new logger
func NewLogger() logger.Logger {
	logger := logrus.New()
	return &stderrLogger{
		internal: logrus.NewEntry(logger),
	}
}

func (l *stderrLogger) Debug(message string) {
	l.internal.Debug(message)
}

func (l *stderrLogger) Info(message string) {
	l.internal.Info(message)
}

func (l *stderrLogger) Warning(message string) {
	l.internal.Warning(message)
}

func (l *stderrLogger) Fatal(message string) {
	l.internal.Fatal(message)
}

func (l *stderrLogger) WithFields(f logger.Fields) logger.Logger {
	return &stderrLogger{
		internal: logrus.WithFields(logrus.Fields(f)),
	}
}
