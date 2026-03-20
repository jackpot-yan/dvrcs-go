package logger

import (
	"fmt"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/sirupsen/logrus"
)

const (
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorGreen  = "\033[32m"
	colorGray   = "\033[37m"
	colorReset  = "\033[0m"
	colorBold   = "\033[1m"
)

type CustomFormatter struct{}

func levelColor(level logrus.Level) string {
	switch level {
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		return colorRed
	case logrus.WarnLevel:
		return colorYellow
	case logrus.InfoLevel:
		return colorGreen
	case logrus.DebugLevel:
		return colorBlue
	default:
		return colorGray
	}
}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	_, file, line, ok := runtime.Caller(7)
	caller := "unknown"
	if ok {
		caller = fmt.Sprintf("%s:%d", filepath.Base(file), line)
	}
	color := levelColor(entry.Level)
	level := fmt.Sprintf("%-5s", entry.Level.String())
	log := fmt.Sprintf("%s[%s]%s %s%s%s  %s%-20s%s  %s\n",
		colorGray, entry.Time.Format("2006-01-02 15:04:05"), colorReset,
		colorBold+color, level, colorReset,
		colorGray, caller, colorReset,
		entry.Message)
	if len(entry.Data) > 0 {
		log = log[:len(log)-1]
		for k, v := range entry.Data {
			log += fmt.Sprintf("  %s%s%s=%v", colorBlue, k, colorReset, v)
		}
		log += "\n"
	}
	return []byte(log), nil
}

var (
	instance *logrus.Logger
	once     sync.Once
)

func getInstance() *logrus.Logger {
	once.Do(func() {
		instance = logrus.New()
		instance.SetFormatter(&CustomFormatter{})
		instance.SetLevel(logrus.DebugLevel)
	})
	return instance
}

func Debug(args ...interface{}) { getInstance().Debug(args...) }
func Info(args ...interface{})  { getInstance().Info(args...) }
func Warn(args ...interface{})  { getInstance().Warn(args...) }
func Error(args ...interface{}) { getInstance().Error(args...) }
func Fatal(args ...interface{}) { getInstance().Fatal(args...) }

func DebugF(format string, args ...interface{}) { getInstance().Debugf(format, args...) }
func Infof(format string, args ...interface{})  { getInstance().Infof(format, args...) }
func Warnf(format string, args ...interface{})  { getInstance().Warnf(format, args...) }
func Errorf(format string, args ...interface{}) { getInstance().Errorf(format, args...) }
func Fatalf(format string, args ...interface{}) { getInstance().Fatalf(format, args...) }

func WithField(key string, value interface{}) *logrus.Entry {
	return getInstance().WithField(key, value)
}

func WithFields(fields logrus.Fields) *logrus.Entry {
	return getInstance().WithFields(fields)
}
