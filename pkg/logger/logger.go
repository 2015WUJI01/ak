package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

type LogStyle string

const (
	ConsoleStyle LogStyle = "console"
	JsonStyle    LogStyle = "json"
)

type LogLevel string

const (
	DebugLevel LogLevel = "debug"
	InfoLevel  LogLevel = "info"
	WarnLevel  LogLevel = "warn"
	ErrorLevel LogLevel = "error"
	PanicLevel LogLevel = "panic"
	FatalLevel LogLevel = "fatal"
)

var l = DefaultLogger()

func DefaultLogger() *zap.Logger {
	cfg := DefaultEncoderConfig()
	return NewLogger(LogConfig{
		Style:         ConsoleStyle,
		Level:         DebugLevel,
		WS:            os.Stdout,
		EncoderConfig: &cfg,
		Hooks:         nil,
	})
}

type LogConfig struct {
	Style         LogStyle
	Level         LogLevel
	WS            zapcore.WriteSyncer
	EncoderConfig *zapcore.EncoderConfig
	Hooks         []func(zapcore.Entry) error
}

func NewLogger(config LogConfig) *zap.Logger {

	if config.EncoderConfig == nil {
		cfg := DefaultEncoderConfig()
		config.EncoderConfig = &cfg
	}

	if config.WS == nil {
		config.WS = os.Stdout
	}

	// 初始化 core
	core := zapcore.NewCore(
		NewEncoder(config.Style, *config.EncoderConfig),
		// os.Stdout, // 日志写入介质
		config.WS,
		logLevelToZapLevel(config.Level),
	)

	// 初始化 logger
	return zap.New(core,
		zap.AddCaller(),      // 调用文件和行号，内部使用 runtime.Caller
		zap.AddCallerSkip(1), // 封装了一层，调用文件去除一层 (runtime.Caller(1))
		zap.Hooks(config.Hooks...),
		// zap.AddStacktrace(zap.ErrorLevel), // Error 时才会显示 stacktrace
	)
}

func SetLogger(config LogConfig) {

	if config.EncoderConfig == nil {
		cfg := DefaultEncoderConfig()
		config.EncoderConfig = &cfg
	}

	if config.WS == nil {
		config.WS = os.Stdout
	}

	// 初始化 core
	core := zapcore.NewCore(
		NewEncoder(config.Style, *config.EncoderConfig),
		// os.Stdout, // 日志写入介质
		config.WS,
		logLevelToZapLevel(config.Level),
	)

	// 初始化 logger
	l = zap.New(core,
		zap.AddCaller(),      // 调用文件和行号，内部使用 runtime.Caller
		zap.AddCallerSkip(1), // 封装了一层，调用文件去除一层 (runtime.Caller(1))
		zap.Hooks(config.Hooks...),
		// zap.AddStacktrace(zap.ErrorLevel), // Error 时才会显示 stacktrace
	)
}

// NewEncoder Zap 日志编码器
func NewEncoder(style LogStyle, cfg zapcore.EncoderConfig) zapcore.Encoder {
	switch style {
	case ConsoleStyle:
		cfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		return zapcore.NewConsoleEncoder(cfg)
	case JsonStyle:
		return zapcore.NewJSONEncoder(cfg)
	default:
		return zapcore.NewConsoleEncoder(cfg)
	}
}

func DefaultEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "Level",
		NameKey:       "logger",
		CallerKey:     "file", // "caller"   代码调用，如 paginator/paginator.go:148
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,        // 每行日志的结尾添加 "\n"
		EncodeLevel:   zapcore.CapitalColorLevelEncoder, // 日志级别名称大小写，如 ERROR、INFO
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) { // 自定义友好的时间格式
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		}, // 时间格式，我们自定义为 2006-01-02 15:04:05.000
		EncodeDuration: zapcore.MillisDurationEncoder, // 执行时间，以秒为单位
		EncodeCaller:   zapcore.ShortCallerEncoder,    // Caller 短格式，如：types/converter.go:17，长格式为绝对路径
	}
}

// 转换 zap 日志级别
func logLevelToZapLevel(level LogLevel) zapcore.Level {
	switch level {
	case DebugLevel:
		return zapcore.DebugLevel
	case InfoLevel:
		return zapcore.InfoLevel
	case WarnLevel:
		return zapcore.WarnLevel
	case ErrorLevel:
		return zapcore.ErrorLevel
	case PanicLevel:
		return zapcore.PanicLevel
	case FatalLevel:
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

func Sync() { _ = l.Sync() }

func Debug(args ...interface{}) { l.Sugar().Debug(args) }
func Debugw(msg string, keysAndValues ...interface{}) {
	l.Sugar().Debugw(msg, keysAndValues...)
}
func Debugf(template string, args ...interface{}) {
	l.Sugar().Debugf(template, args...)
}

func Info(args ...interface{}) { l.Sugar().Info(args) }
func Infow(msg string, keysAndValues ...interface{}) {
	l.Sugar().Infow(msg, keysAndValues...)
}
func Infof(template string, args ...interface{}) {
	l.Sugar().Infof(template, args...)
}

func Warn(args ...interface{}) { l.Sugar().Warn(args) }
func Warnw(msg string, keysAndValues ...interface{}) {
	l.Sugar().Warnw(msg, keysAndValues...)
}
func Warnf(template string, args ...interface{}) {
	l.Sugar().Warnf(template, args...)
}

func Error(args ...interface{}) { l.Sugar().Error(args) }
func Errorw(msg string, keysAndValues ...interface{}) {
	l.Sugar().Errorw(msg, keysAndValues...)
}
func Errorf(template string, args ...interface{}) {
	l.Sugar().Errorf(template, args...)
}

func Panic(args ...interface{}) { l.Sugar().Panic(args) }
func Panicw(msg string, keysAndValues ...interface{}) {
	l.Sugar().Panicw(msg, keysAndValues...)
}
func Panicf(template string, args ...interface{}) {
	l.Sugar().Panicf(template, args...)
}

func Fatal(args ...interface{}) { l.Sugar().Fatal(args) }
func Fatalw(msg string, keysAndValues ...interface{}) {
	l.Sugar().Fatalw(msg, keysAndValues...)
}
func Fatalf(template string, args ...interface{}) {
	l.Sugar().Fatalf(template, args...)
}
