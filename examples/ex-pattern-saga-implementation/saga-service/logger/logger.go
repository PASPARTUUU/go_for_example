package logger

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// New -
func New() *logrus.Logger {
	loger := logrus.New()

	// default logging level
	loger.SetLevel(logrus.InfoLevel)

	return loger
}

// SetLogLevel - sets level for the logger.
func SetLogLevel(l *logrus.Logger, level string) error {
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		return errors.Wrap(err, "failed to parse a log level")
	}
	l.SetLevel(lvl)
	return nil
}
