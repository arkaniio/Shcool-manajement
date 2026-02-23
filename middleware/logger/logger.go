package logger

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

// Init initializes the logger with production best practices
func Init() {
	// Get working directory for log files
	wd, _ := os.Getwd()
	logDir := filepath.Join(wd, "logs")

	// Create logs directory if it doesn't exist
	if err := os.MkdirAll(logDir, 0755); err != nil {
		panic(err)
	}

	// Configure encoder for human-friendly output
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Define log levels
	level := zapcore.InfoLevel

	// Create encoder
	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	// Create file writer for persistent logs
	file, err := os.OpenFile(
		filepath.Join(logDir, "app.log"),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		panic(err)
	}

	// Create multi-write: console + file
	writeSyncer := zapcore.AddSync(file)

	// Create core with both console and file output
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level),   // Console
		zapcore.NewCore(encoder, writeSyncer, level),                  // File
	)

	// Create logger with caller info
	Log = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
}

// InitProduction initializes logger optimized for production
func InitProduction() {
	// Get working directory for log files
	wd, _ := os.Getwd()
	logDir := filepath.Join(wd, "logs")

	// Create logs directory if it doesn't exist
	if err := os.MkdirAll(logDir, 0755); err != nil {
		panic(err)
	}

	// Production encoder config (JSON format)
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Production level
	level := zapcore.InfoLevel

	// JSON encoder for production
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	// File for app logs
	appFile, err := os.OpenFile(
		filepath.Join(logDir, "app.json.log"),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		panic(err)
	}

	// Separate file for error logs
	errorFile, err := os.OpenFile(
		filepath.Join(logDir, "error.json.log"),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		panic(err)
	}

	// Create cores for different log levels
	core := zapcore.NewTee(
		// Info and above to app file
		zapcore.NewCore(encoder, zapcore.AddSync(appFile), level),
		// Error and above to error file
		zapcore.NewCore(encoder, zapcore.AddSync(errorFile), zapcore.ErrorLevel),
	)

	Log = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
}

// Sync flushes any buffered log entries
func Sync() {
	if Log != nil {
		_ = Log.Sync()
	}
}

// WithFields returns a logger with additional fields
func WithFields(fields ...zapcore.Field) *zap.Logger {
	return Log.With(fields...)
}

