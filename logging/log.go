package logging

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	InfoPrefix    = "[LOG.KX][INFO]"
	ErrorPrefix   = "[LOG.KX][ERROR]"
	WarningPrefix = "[LOG.KX][WARNING]"
	DebugPrefix   = "[LOG.KX][DEBUG]"
	FatalPrefix   = "[LOG.KX][FATAL]"
)

var logger *zap.SugaredLogger
var logError error

var logSync sync.Once

func configure(path string, maxSize int, maxBackup int) zapcore.WriteSyncer {
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   path,
		MaxSize:    maxSize,
		MaxBackups: maxBackup,
		Compress:   false,
	})
	return zapcore.NewMultiWriteSyncer(
		zapcore.AddSync(os.Stdout),
		zapcore.AddSync(w),
	)
}

func GetLogger(conf LogConfig) error {
	logSync.Do(func() {
		encoderCfg := getApplicationEncoderConfig()

		var encoder zapcore.Encoder

		encoder = zapcore.NewJSONEncoder(encoderCfg)

		if conf.Formatter == "console" {
			encoder = zapcore.NewConsoleEncoder(encoderCfg)
		}

		var level zap.AtomicLevel
		level, logError = zap.ParseAtomicLevel(conf.Level)
		if logError != nil {
			return
		}

		core := zapcore.NewCore(encoder, configure(conf.Path, conf.MaxSize, conf.MaxBackups), level)
		loggerZap := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
		logger = loggerZap.Sugar()
	})

	return logError
}

func getApplicationEncoderConfig() zapcore.EncoderConfig {
	var en = zap.NewProductionEncoderConfig()
	en.LevelKey = "level"
	en.CallerKey = "caller"
	en.TimeKey = "timestamp"
	en.MessageKey = "message"
	en.EncodeDuration = zapcore.NanosDurationEncoder
	en.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.UTC().Format("2006-01-02T15:04:05.000Z"))
	}
	return en
}

func getMessage(template string, fmtArgs []interface{}) string {
	if len(fmtArgs) == 0 {
		return template
	}

	if template != "" {
		return fmt.Sprintf(template, fmtArgs...)
	}

	if len(fmtArgs) == 1 {
		if str, ok := fmtArgs[0].(string); ok {
			return str
		}
	}
	return fmt.Sprint(fmtArgs...)
}

func Infof(context context.Context, format string, args ...interface{}) {
	if logger == nil {
		fmt.Printf(format+"\n", args...)
		return
	}

	pairs := GetLogFieldsFromContext(context)
	msg := fmt.Sprintf("%s %s", InfoPrefix, getMessage(format, args))

	logger.Infow(msg, pairs...)
}

func InfofNw(format string, args ...interface{}) {
	if logger == nil {
		fmt.Printf(format+"\n", args...)
		return
	}

	msg := fmt.Sprintf("%s %s", InfoPrefix, getMessage(format, args))

	logger.Info(msg)
}

func Debugf(context context.Context, format string, args ...interface{}) {
	if logger == nil {
		fmt.Printf(format+"\n", args...)
		return
	}

	pairs := GetLogFieldsFromContext(context)
	msg := fmt.Sprintf("%s %s", DebugPrefix, getMessage(format, args))

	logger.Debugw(msg, pairs...)
}

func Fatalf(format string, args ...interface{}) {
	if logger == nil {
		fmt.Printf(format+"\n", args...)
		return
	}

	msg := fmt.Sprintf("%s %s", FatalPrefix, getMessage(format, args))

	logger.Fatalf(msg)
}

func DebugfNw(format string, args ...interface{}) {
	if logger == nil {
		fmt.Printf(format+"\n", args...)
		return
	}

	msg := fmt.Sprintf("%s %s", DebugPrefix, getMessage(format, args))

	logger.Debug(msg)
}

func Warn(context context.Context, msg string) {
	if logger == nil {
		fmt.Println(msg)
		return
	}

	pairs := GetLogFieldsFromContext(context)
	logger.Warnw(fmt.Sprintf("%s %s", WarningPrefix, msg), pairs...)
}

func Warnf(context context.Context, format string, args ...interface{}) {
	if logger == nil {
		fmt.Printf(format+"\n", args...)
		return
	}

	pairs := GetLogFieldsFromContext(context)
	msg := fmt.Sprintf("%s %s", WarningPrefix, getMessage(format, args))

	logger.Warnw(msg, pairs...)
}

func WarnfNw(format string, args ...interface{}) {
	if logger == nil {
		fmt.Printf(format+"\n", args...)
		return
	}

	msg := fmt.Sprintf("%s %s", WarningPrefix, getMessage(format, args))
	logger.Warn(msg)
}

func Errorf(context context.Context, format string, args ...interface{}) {
	if logger == nil {
		fmt.Printf(format+"\n", args...)
		return
	}

	pairs := GetLogFieldsFromContext(context)
	msg := fmt.Sprintf("%s %s", ErrorPrefix, getMessage(format, args))

	logger.Warnw(msg, pairs...)
}

func GetLogFieldsFromContext(context context.Context) []interface{} {
	var pairs []interface{}
	if context != nil {
		pairs = []interface{}{}
	}
	return pairs
}

func Log(context context.Context, code int, msg string) {
	if logger == nil {
		fmt.Println(msg)
		return
	}

	pairs := GetLogFieldsFromContext(context)

	logger.Warnw(fmt.Sprintf("%s %s", WarningPrefix, msg), pairs...)
}

func Logf(context context.Context, code int, format string, args ...interface{}) {
	if logger == nil {
		fmt.Printf(format+"\n", args...)
		return
	}

	pairs := GetLogFieldsFromContext(context)
	msg := fmt.Sprintf("%s %s", WarningPrefix, getMessage(format, args))
	logger.Warnw(msg, pairs...)
}

func LogfNw(code int, format string, args ...interface{}) {
	if logger == nil {
		fmt.Printf(format+"\n", args...)
		return
	}
	msg := fmt.Sprintf("%s %s", WarningPrefix, getMessage(format, args))
	logger.Warn(msg)
}
