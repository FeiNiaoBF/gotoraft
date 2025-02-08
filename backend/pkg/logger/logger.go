// pkg/logger/logger.go
// 日志管理包，提供全局日志记录功能
// 使用zap（go.uber.org/zap）

package logger

import (
	"fmt"
	"gotoraft/config"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	logger *zap.Logger
	sugar  *zap.SugaredLogger
)

// InitLogger 初始化日志系统
func InitLogger() error {
	cfg := config.GetLogConfig()
	if cfg == nil {
		return fmt.Errorf("failed to get logger config")
	}

	// 日志级别
	var level zapcore.Level
	switch cfg.Level {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}

	// 设置编码器配置
	// 用来指定日志的格式
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:          "time",
		LevelKey:         "level",
		NameKey:          "logger",
		CallerKey:        "caller",
		MessageKey:       "msg",
		StacktraceKey:    "stacktrace",
		LineEnding:       zapcore.DefaultLineEnding,
		EncodeLevel:      zapcore.CapitalLevelEncoder,
		EncodeTime:       customTimeEncoder(cfg.TimeFormat),
		EncodeDuration:   zapcore.StringDurationEncoder,
		EncodeCaller:     zapcore.ShortCallerEncoder,
		EncodeName:       zapcore.FullNameEncoder,
		ConsoleSeparator: " | ",
	}

	// 创建编码器
	var encoder zapcore.Encoder
	if cfg.Format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// 设置输出
	var cores []zapcore.Core

	// 控制台输出
	if cfg.Output == "console" || cfg.Output == "both" {
		consoleCore := zapcore.NewCore(
			encoder,
			zapcore.AddSync(os.Stdout),
			level,
		)
		cores = append(cores, consoleCore)
	}

	// 文件输出
	if cfg.Output == "file" || cfg.Output == "both" {
		// 根据当前时间生成日志文件路径
		now := time.Now()
		timeBasedPath := filepath.Join(
			filepath.Dir(cfg.Filename),
			fmt.Sprintf("%d/%02d/%02d", now.Year(), now.Month(), now.Day()),
		)

		// 确保日志目录存在
		if err := os.MkdirAll(timeBasedPath, 0755); err != nil {
			return fmt.Errorf("failed to create log directory: %v", err)
		}

		// 生成带时间戳的日志文件名
		logFileName := filepath.Join(
			timeBasedPath,
			fmt.Sprintf("%s_%s.log",
				filepath.Base(cfg.Filename[:len(cfg.Filename)-len(filepath.Ext(cfg.Filename))]), // 移除.log后缀
				now.Format("15-04"), // 添加时分
			),
		)

		// 配置日志轮转
		writer := &lumberjack.Logger{
			Filename:   logFileName,
			MaxSize:    cfg.MaxSize, // MB
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge, // days
			Compress:   cfg.Compress,
		}

		fileCore := zapcore.NewCore(
			encoder,
			zapcore.AddSync(writer),
			level,
		)
		cores = append(cores, fileCore)
	}

	// 创建核心
	core := zapcore.NewTee(cores...)

	// 创建logger
	logger = zap.New(core,
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)

	sugar = logger.Sugar()
	return nil
}

// customTimeEncoder 自定义时间编码器
func customTimeEncoder(format string) zapcore.TimeEncoder {
	if format == "" {
		format = "2006-01-02 15:04:05"
	}
	return func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(format))
	}
}

// Debug 输出Debug级别日志
func Debug(args ...interface{}) {
	sugar.Debug(args...)
}

// Debugf 输出Debug级别日志（格式化）
func Debugf(template string, args ...interface{}) {
	sugar.Debugf(template, args...)
}

// Info 输出Info级别日志
func Info(args ...interface{}) {
	sugar.Info(args...)
}

// Infof 输出Info级别日志（格式化）
func Infof(template string, args ...interface{}) {
	sugar.Infof(template, args...)
}

// Warn 输出Warn级别日志
func Warn(args ...interface{}) {
	sugar.Warn(args...)
}

// Warnf 输出Warn级别日志（格式化）
func Warnf(template string, args ...interface{}) {
	sugar.Warnf(template, args...)
}

// Error 输出Error级别日志
func Error(args ...interface{}) {
	sugar.Error(args...)
}

// Errorf 输出Error级别日志（格式化）
func Errorf(template string, args ...interface{}) {
	sugar.Errorf(template, args...)
}

// Fatal 输出Fatal级别日志
func Fatal(args ...interface{}) {
	sugar.Fatal(args...)
}

// Fatalf 输出Fatal级别日志（格式化）
func Fatalf(template string, args ...interface{}) {
	sugar.Fatalf(template, args...)
}

// Sync 同步日志缓冲区
func Sync() error {
	if sugar != nil {
		return sugar.Sync()
	}
	return nil
}
