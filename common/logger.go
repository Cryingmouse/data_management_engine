package common

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/gorm/logger"
)

var (
	AuditLogger *log.Logger
	Logger      *log.Logger
	DBLogger    *LogrusLogger
)

func SetupLoggers() {
	// Create a logger for audit logging
	AuditLogger = log.New()
	AuditLogger.SetFormatter(&log.JSONFormatter{})
	// 设置logrus的输出为lumberjack日志轮转器
	AuditLogger.SetOutput(&lumberjack.Logger{
		Filename:   Config.Logger.AuditLogFile,
		MaxSize:    10, // megabytes
		MaxBackups: 3,
		MaxAge:     28,   //days
		Compress:   true, // disabled by default
	})

	// Create a logger for general logging (JSON format)
	Logger = log.New()
	Logger.SetFormatter(&log.JSONFormatter{})
	// 打开日志文件
	Logger.SetOutput(&lumberjack.Logger{
		Filename:   Config.Logger.LogFile,
		MaxSize:    10, // megabytes
		MaxBackups: 3,
		MaxAge:     28,   //days
		Compress:   true, // disabled by default
	})

	// Set the log levels for both loggers independently
	setLoggerLevel()

	DBLogger = &LogrusLogger{log: Logger}
}

func setLoggerLevel() {
	switch Config.Logger.LogLevel {
	case "error":
		Logger.SetLevel(log.ErrorLevel)
	case "info":
		Logger.SetLevel(log.InfoLevel)
	case "debug":
		Logger.SetLevel(log.DebugLevel)
	case "trace":
		Logger.SetLevel(log.TraceLevel)
	default:
		Logger.SetLevel(log.InfoLevel)
	}
}

type LogrusLogger struct {
	log *log.Logger
}

// LogMode sets the log mode for the logger
func (l *LogrusLogger) LogMode(mode logger.LogLevel) logger.Interface {
	return l
}

// Info logs an info message
func (l *LogrusLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	l.log.Infof(msg, data...)
}

// Warn logs a warning message
func (l *LogrusLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	l.log.Warnf(msg, data...)
}

// Error logs an error message
func (l *LogrusLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	l.log.Errorf(msg, data...)
}

// Trace logs an SQL statement and its execution time
func (l *LogrusLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	sql, rows := fc()
	l.log.WithFields(log.Fields{
		"context":  ctx,
		"rows":     rows,
		"duration": time.Since(begin),
		"error":    err,
	}).Tracef("SQL: %s", sql)
}
