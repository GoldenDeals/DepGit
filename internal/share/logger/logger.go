package logger

import (
	"github.com/sirupsen/logrus"
)

type customLogger struct {
	defaultField string
	formatter    logrus.Formatter
}

func (l customLogger) Format(entry *logrus.Entry) ([]byte, error) {
	entry.Data["src"] = l.defaultField
	return l.formatter.Format(entry)
}

func New(name string, logLevel ...logrus.Level) *logrus.Logger {
	log := logrus.New()
	log.SetFormatter(customLogger{
		defaultField: name,
		formatter:    logrus.StandardLogger().Formatter,
	})
	if len(logLevel) > 0 {
		log.SetLevel(logLevel[0])
	} else {
		log.SetLevel(logrus.DebugLevel)
	}

	return log
}
