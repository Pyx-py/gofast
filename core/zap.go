package core

import (
	"fmt"
	"os"
	"time"

	"github.com/Pyx-py/gofast/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogConfig struct {
	LogPath      string // 日志文件路径，为空不启用默认日志
	LogLevel     string // 日志等级
	LogName      string // 日志文件的前缀名，也是最新的日志文件的软链接名称
	EncoderLevel string // 日志编码器的类型，默认为小写不带色彩
	TextType     string // 日志类型，json或者是console
	Day          int    // 日志保存几天
	LogInConsole bool   // 是否在控制台同时打印日志
}

func (lc LogConfig) initConfig() (nlc LogConfig, err error) {
	if lc.LogPath == "" {
		return nlc, fmt.Errorf("logPath can't be empty")
	}
	if lc.LogLevel == "" {
		nlc.LogLevel = "info"
	} else {
		nlc.LogLevel = lc.LogLevel
	}
	if lc.LogName == "" {
		nlc.LogName = "gofast"
	} else {
		nlc.LogName = lc.LogName
	}
	if lc.EncoderLevel == "" {
		nlc.EncoderLevel = "LowercaseLevelEncoder"
	} else {
		nlc.EncoderLevel = lc.EncoderLevel
	}
	if lc.TextType == "" {
		nlc.TextType = "console"
	} else {
		nlc.TextType = lc.TextType
	}
	return nlc, nil
}

var level zapcore.Level

func (lc LogConfig) Zap() (logger *zap.Logger) {
	nlc, err := lc.initConfig()
	if err != nil {
		panic(err)
	}
	if ok, _ := utils.PathExists(nlc.LogPath); !ok {
		fmt.Printf("create %v directory\n", nlc.LogPath)
		_ = os.Mkdir(nlc.LogPath, os.ModePerm)
	}

	switch nlc.LogLevel {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	case "dpanic":
		level = zap.DPanicLevel
	case "panic":
		level = zap.PanicLevel
	case "fatal":
		level = zap.FatalLevel
	default:
		level = zap.InfoLevel
	}

	if level == zap.DebugLevel || level == zap.ErrorLevel {
		logger = zap.New(lc.getEncoderCore(), zap.AddStacktrace(level))
	} else {
		logger = zap.New(lc.getEncoderCore())
	}
	logger = logger.WithOptions(zap.AddCaller())
	return logger
}

func (lc LogConfig) getEncoderConfig() (config zapcore.EncoderConfig) {
	config = zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     CustomTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}

	switch lc.EncoderLevel {
	case "LowercaseLevelEncoder":
		config.EncodeLevel = zapcore.LowercaseLevelEncoder
	case "LowercaseColorLevelEncoder":
		config.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	case "CapitalLevelEncoder":
		config.EncodeLevel = zapcore.CapitalLevelEncoder
	case "CapitalColorLevelEncoder":
		config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	default:
		config.EncodeLevel = zapcore.LowercaseLevelEncoder

	}
	return config
}

func (lc LogConfig) getEncoder() zapcore.Encoder {
	if lc.TextType == "json" {
		return zapcore.NewJSONEncoder(lc.getEncoderConfig())
	}
	return zapcore.NewConsoleEncoder(lc.getEncoderConfig())
}

func (lc LogConfig) getEncoderCore() (core zapcore.Core) {
	writer, err := utils.GetWriteSyncer(lc.LogName, lc.LogPath, lc.Day, lc.LogInConsole)
	if err != nil {
		fmt.Printf("get write syncer failed err:%v", err.Error())
		return
	}
	return zapcore.NewCore(lc.getEncoder(), writer, level)
}

func CustomTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006/01/02-15:04:05.000"))
}
