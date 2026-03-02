package logger

import "github.com/sirupsen/logrus"

func Init(cfg *LoggerConfig) {

	switch cfg.Level {
	case "info":
		{
			logrus.SetLevel(logrus.InfoLevel)
		}
	case "debug":
		{
			logrus.SetLevel(logrus.DebugLevel)
		}
	case "warn":
		{
			logrus.SetLevel(logrus.WarnLevel)
		}
	case "error":
		{
			logrus.SetLevel(logrus.ErrorLevel)
		}
	default:
		{
			logrus.SetLevel(logrus.ErrorLevel)
		}
	}

	if cfg.JSONFormatter {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}

	logrus.SetReportCaller(cfg.ShowMethod)

}

func Debug(msg ...interface{}) {
	logrus.Debug(msg...)
}

func Debugf(format string, args ...interface{}) {
	logrus.Debugf(format, args...)
}

func Info(msg ...interface{}) {
	logrus.Info(msg...)
}

func Infof(format string, args ...interface{}) {
	logrus.Infof(format, args...)
}

func Warn(msg ...interface{}) {
	logrus.Warn(msg...)
}

func Warnf(format string, args ...interface{}) {
	logrus.Warnf(format, args...)
}

func Error(msg ...interface{}) {
	logrus.Error(msg...)
}

func Errorf(format string, args ...interface{}) {
	logrus.Errorf(format, args...)
}

func Fatal(msg ...interface{}) {
	logrus.Fatal(msg...)
}

func Fatalf(format string, msg ...interface{}) {
	logrus.Fatalf(format, msg...)
}
