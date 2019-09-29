package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/mikloslorinczi/pager/version"

	"github.com/sirupsen/logrus"
)

// Logger wrapper for logrus.FieldLogger
type Logger struct {
	logrus.FieldLogger
}

var defaultLogger Logger

func init() {
	baseLogger := *logrus.New()

	// Set log level
	logLevel := "info"
	if envLvl := os.Getenv("LOG_LEVEL"); envLvl != "" {
		logLevel = envLvl
	}
	lvl, err := logrus.ParseLevel(logLevel)
	if err != nil {
		logrus.Fatal("Error when parsing LOG_LEVEL environment variable")
	}
	baseLogger.SetLevel(lvl)

	// Set time format
	baseLogger.Formatter = &logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339Nano,
	}

	// Set output
	baseLogger.Out = os.Stdout

	// Set default fileds
	defaultFields := logrus.Fields{
		"App": "Pager",
		"Ver": version.Get(),
	}

	defaultLogger = Logger{baseLogger.WithFields(defaultFields)}

	defaultLogger.Debug("Logger initialized")
}

// NewLogger return defaultLogger
func NewLogger(component string) Logger {
	return defaultLogger.setComponent(component)
}

func (l Logger) setComponent(comp string) Logger {
	return Logger{l.FieldLogger.WithField("component", comp)}
}

// WithField is a wrapper to return a Logger
func (l Logger) WithField(s string, i interface{}) Logger {
	return Logger{l.FieldLogger.WithField(s, i)}
}

// WithFields is a wrapper to return a Logger
func (l Logger) WithFields(m map[string]interface{}) Logger {
	return Logger{l.FieldLogger.WithFields(m)}
}

// ErrorfAndStacktrace ...
func (l Logger) ErrorfAndStacktrace(err error, s string, args ...interface{}) {
	l.WithField("error", err.Error()).WithField("stacktrace", fmt.Sprintf("%+v", err)).Errorf(s, args...)
}

// ErrorAndStacktrace ...
func (l Logger) ErrorAndStacktrace(err error, args ...interface{}) {
	l.WithField("error", err.Error()).WithField("stacktrace", fmt.Sprintf("%+v", err)).Error(args...)
}

// WithError calls logrus.WithError and wraps the error into errors.WithStack
func (l Logger) WithError(err error) Logger {
	return l.WithField("error", err.Error()).WithField("stacktrace", fmt.Sprintf("%+v", err))
}
