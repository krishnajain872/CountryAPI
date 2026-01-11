package logger

import (
	"fmt"
	"os"
	"runtime"
	"github.com/krishnajain872/country-search-api-cache/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// ==============================
// Logger wrapper
// ==============================
type Logger struct {
	zap      *zap.Logger
	env      config.LogEnvType
	severity zapcore.Level
	modes    map[config.LogModeType]bool
}

// ==============================
// Constructor
// ==============================
func NewLogger(cfg *config.LoggerConfig) *Logger {
	sevLevel := parseSeverity(cfg.Severity)

	encoderCfg := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	var cores []zapcore.Core
	for _, mode := range cfg.Mode {
		switch mode {
		case config.FileMode:
			lumberjackLogger := &lumberjack.Logger{
				Filename:   cfg.FilePath,
				MaxSize:    cfg.MaxSizeMB,
				MaxBackups: cfg.MaxBackups,
				MaxAge:     cfg.MaxAgeDays,
				Compress:   true,
			}
			fileCore := zapcore.NewCore(
				zapcore.NewJSONEncoder(encoderCfg),
				zapcore.AddSync(lumberjackLogger),
				sevLevel,
			)
			cores = append(cores, fileCore)

		case config.ConsoleMode:
			consoleCore := zapcore.NewCore(
				zapcore.NewConsoleEncoder(encoderCfg),
				zapcore.Lock(os.Stdout),
				sevLevel,
			)
			cores = append(cores, consoleCore)
		}
	}

	options := []zap.Option{
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	}
	if cfg.Environment == config.Development {
		options = append(options, zap.Development())
	}

	return &Logger{
		zap:      zap.New(zapcore.NewTee(cores...), options...),
		env:      cfg.Environment,
		severity: sevLevel,
		modes:    modesMap(cfg.Mode),
	}
}

// ==============================
// Singleton helper
// ==============================
var logInstance *Logger

func GetLogger() *Logger {
	if logInstance == nil {
		cfg := config.Get().Logger
		logInstance = NewLogger(cfg)
	}
	return logInstance
}

// ==============================
// Logging methods
// ==============================
func (l *Logger) Debug(msg string, fields ...zap.Field) {
	if l.shouldLog(zapcore.DebugLevel) {
		l.zap.Debug(msg, append(fields, zap.Int64("gid", int64(getGID())))...)
	}
}

func (l *Logger) Info(msg string, fields ...zap.Field) {
	if l.shouldLog(zapcore.InfoLevel) {
		l.zap.Info(msg, append(fields, zap.Int64("gid", int64(getGID())))...)
	}
}

func (l *Logger) Warn(msg string, fields ...zap.Field) {
	if l.shouldLog(zapcore.WarnLevel) {
		l.zap.Warn(msg, append(fields, zap.Int64("gid", int64(getGID())))...)
	}
}

func (l *Logger) Error(msg string, fields ...zap.Field) {
	if l.shouldLog(zapcore.ErrorLevel) {
		l.zap.Error(msg, append(fields, zap.Int64("gid", int64(getGID())))...)
	}
}

func (l *Logger) DPanic(msg string, fields ...zap.Field) {
	if l.shouldLog(zapcore.DPanicLevel) {
		l.zap.DPanic(msg, append(fields, zap.Int64("gid", int64(getGID())))...)
	}
}

func (l *Logger) Panic(msg string, fields ...zap.Field) {
	if l.shouldLog(zapcore.PanicLevel) {
		l.zap.Panic(msg, append(fields, zap.Int64("gid", int64(getGID())))...)
	}
}

func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	if l.shouldLog(zapcore.FatalLevel) {
		l.zap.Fatal(msg, append(fields, zap.Int64("gid", int64(getGID())))...)
	}
}

// Sync flushes all logs
func (l *Logger) Sync() {
	_ = l.zap.Sync()
}

// ==============================
// Helpers
// ==============================

func parseSeverity(sev config.LogSeverity) zapcore.Level {
	switch sev {
	case config.DebugSeverity:
		return zapcore.DebugLevel
	case config.InfoSeverity:
		return zapcore.InfoLevel
	case config.WarnSeverity:
		return zapcore.WarnLevel
	case config.ErrorSeverity:
		return zapcore.ErrorLevel
	case config.CriticalSeverity:
		return zapcore.DPanicLevel
	default:
		return zapcore.InfoLevel
	}
}

func modesMap(modes []config.LogModeType) map[config.LogModeType]bool {
	m := make(map[config.LogModeType]bool)
	for _, v := range modes {
		m[v] = true
	}
	return m
}

func (l *Logger) shouldLog(level zapcore.Level) bool {
	if l.env == config.Development {
		return true
	}
	return level >= l.severity
}

// getGID returns current goroutine ID
func getGID() uint64 {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	var gid uint64
	fmt.Sscanf(string(buf[:n]), "goroutine %d ", &gid)
	return gid
}

// ==============================
// Zap field helpers (for structured logging)
// ==============================
func String(key, val string) zap.Field {
	return zap.String(key, val)
}

func Any(key string, val interface{}) zap.Field {
	return zap.Any(key, val)
}

func Int(key string, val int) zap.Field {
	return zap.Int(key, val)
}
