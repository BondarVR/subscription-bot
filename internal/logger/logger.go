package logger

import (
	graylog "github.com/gemnasium/logrus-graylog-hook"
	logrus "github.com/sirupsen/logrus"
	"time"
)

type Logger interface {
	Info(args ...interface{})
	Fatal(args ...interface{})
	Infof(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}

type Config struct {
	LogLevel    string
	LogServer   string
	ServiceName string
}

type LogrusLogger struct {
	logrus *logrus.Logger
	entry  *logrus.Entry
}

func (l *LogrusLogger) Info(args ...interface{}) {
	l.entry.Info(args...)
}

func (l *LogrusLogger) Fatal(args ...interface{}) {
	l.entry.Fatal(args...)
}

func (l *LogrusLogger) Infof(format string, args ...interface{}) {
	l.entry.Infof(format, args...)
}

func (l *LogrusLogger) Fatalf(format string, args ...interface{}) {
	l.entry.Fatalf(format, args...)
}

func New(cfg Config) (*LogrusLogger, error) {
	level, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		return nil, err
	}

	logger := &LogrusLogger{
		logrus: logrus.New(),
	}

	logger.logrus.SetLevel(level)

	customFormatter := &logrus.JSONFormatter{
		TimestampFormat: time.Layout,
	}
	logger.logrus.SetFormatter(customFormatter)

	if cfg.LogServer != "" {
		logger.logrus.AddHook(
			graylog.NewGraylogHook(cfg.LogServer, map[string]interface{}{}),
		)
	}

	logger.entry = logger.logrus.WithFields(logrus.Fields{
		"service_name": cfg.ServiceName,
	})

	return logger, nil
}
