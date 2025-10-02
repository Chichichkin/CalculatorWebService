package logger

import (
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"CalculatorWebService/internal/config"
)

var Logger *logrus.Logger

func InitLogger(config config.LoggerConfig) {
	Logger = logrus.New()

	level, err := logrus.ParseLevel(config.Level)
	if err != nil {
		Logger.Warnf("Invalid log level '%s', defaulting to 'info'", config.Level)
		level = logrus.InfoLevel
	}
	Logger.SetLevel(level)

	switch strings.ToLower(config.Format) {
	case "json":
		Logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "timestamp",
				logrus.FieldKeyLevel: "level",
				logrus.FieldKeyMsg:   "message",
			},
		})
	case "text":
		Logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: config.TimeFormat, //"2006-01-02 15:04:05",
		})
	default:
		Logger.Warnf("Invalid log format '%s', defaulting to 'text'", config.Format)
		Logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: config.TimeFormat,
		})
	}

	Logger.SetOutput(os.Stdout)

	Logger = Logger.WithFields(logrus.Fields{
		"service": config.ServerName, // "calculator"
		"version": config.Version,    // "1.0.0"
	}).Logger
}

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()
		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		bodySize := c.Writer.Size()

		entry := Logger.WithFields(logrus.Fields{
			"method":     method,
			"path":       path,
			"query":      raw,
			"status":     statusCode,
			"latency":    latency,
			"client_ip":  clientIP,
			"body_size":  bodySize,
			"user_agent": c.Request.UserAgent(),
		})

		if statusCode >= 500 {
			entry.Error("HTTP request completed with server error")
		} else if statusCode >= 400 {
			entry.Warn("HTTP request completed with client error")
		} else {
			entry.Info("HTTP request completed successfully")
		}
	}
}

func LogInfo(message string, fields ...logrus.Fields) {
	if len(fields) > 0 {
		Logger.WithFields(fields[0]).Info(message)
	} else {
		Logger.Info(message)
	}
}

func LogWarn(message string, fields ...logrus.Fields) {
	if len(fields) > 0 {
		Logger.WithFields(fields[0]).Warn(message)
	} else {
		Logger.Warn(message)
	}
}

func LogError(message string, err error, fields ...logrus.Fields) {
	logFields := logrus.Fields{}
	if err != nil {
		logFields["error"] = err.Error()
	}
	if len(fields) > 0 {
		for k, v := range fields[0] {
			logFields[k] = v
		}
	}
	Logger.WithFields(logFields).Error(message)
}
