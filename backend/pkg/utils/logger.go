package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

// InitLogger تهيئة الـ logger
func InitLogger() {
	Logger = logrus.New()

	// تعيين الـ output
	Logger.SetOutput(os.Stdout)

	// تعيين الـ format (JSON للإنتاج، Text للتطوير)
	if os.Getenv("APP_ENV") == "production" {
		Logger.SetFormatter(&logrus.JSONFormatter{})
	} else {
		Logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
			ForceColors:   true,
		})
	}

	// تعيين الـ log level
	level := os.Getenv("LOG_LEVEL")
	if level == "" {
		level = "info"
	}

	switch level {
	case "debug":
		Logger.SetLevel(logrus.DebugLevel)
	case "info":
		Logger.SetLevel(logrus.InfoLevel)
	case "warn":
		Logger.SetLevel(logrus.WarnLevel)
	case "error":
		Logger.SetLevel(logrus.ErrorLevel)
	default:
		Logger.SetLevel(logrus.InfoLevel)
	}
}

// InitLoggerWithConfig تهيئة الـ logger مع إعدادات محددة
func InitLoggerWithConfig(logLevel, appEnv string) {
	os.Setenv("LOG_LEVEL", logLevel)
	os.Setenv("APP_ENV", appEnv)
	InitLogger()
}

// LogInfo تسجيل معلومات
func LogInfo(message string, fields logrus.Fields) {
	if Logger == nil {
		InitLogger()
	}
	Logger.WithFields(fields).Info(message)
}

// LogError تسجيل خطأ
func LogError(message string, fields logrus.Fields) {
	if Logger == nil {
		InitLogger()
	}
	Logger.WithFields(fields).Error(message)
}

// LogDebug تسجيل معلومات للتصحيح
func LogDebug(message string, fields logrus.Fields) {
	if Logger == nil {
		InitLogger()
	}
	Logger.WithFields(fields).Debug(message)
}

// LogWarning تسجيل تحذير
func LogWarning(message string, fields logrus.Fields) {
	if Logger == nil {
		InitLogger()
	}
	Logger.WithFields(fields).Warning(message)
}
