package main

import (
	"fmt"
	"os"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var _ log.Logger = (*ZapLogger)(nil)

// Logger 配置 zap 日志,将 zap 日志库引入
func Logger() log.Logger {
	//配置 zap 日志库的编码器
	encoder := zapcore.EncoderConfig{
		TimeKey:        "time", // 定义时间标签的名字
		LevelKey:       "level",
		NameKey:        "logger",
		MessageKey:     "msg",
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
	return NewZapLogger(
		encoder,
		zap.NewAtomicLevelAt(zapcore.DebugLevel),
		zap.AddStacktrace(zap.NewAtomicLevelAt(zapcore.ErrorLevel)),
		// zap.AddCaller(),
		// zap.AddCallerSkip(2),
		zap.Development(),
	)
}

// 日志自动切割，采用 lumberjack 实现的
func getLogWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "./log/zap.log",
		MaxSize:    10,    //日志的最大大小（M）
		MaxBackups: 5,     //日志的最大保存数量
		MaxAge:     30,    //日志文件存储最大天数
		Compress:   false, //是否执行压缩
	}
	return zapcore.AddSync(lumberJackLogger)
}

type ZapLogger struct {
	log  *zap.Logger
	Sync func() error
}

// NewZapLogger return a zap logger.
func NewZapLogger(encoder zapcore.EncoderConfig, level zap.AtomicLevel, opts ...zap.Option) *ZapLogger {
	//日志切割
	writeSyncer := getLogWriter()
	//设置日志级别
	level.SetLevel(zap.DebugLevel)
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoder), // 编码器配置
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(writeSyncer)), // 打印到控制台和文件
		level, // 日志级别
	)
	zapLogger := zap.New(core, opts...)
	return &ZapLogger{log: zapLogger, Sync: zapLogger.Sync}
}

// Log 实现log接口
func (l *ZapLogger) Log(level log.Level, keyvals ...interface{}) error {
	if len(keyvals) == 0 || len(keyvals)%2 != 0 {
		l.log.Warn(fmt.Sprint("Keyvalues must appear in pairs: ", keyvals))
		return nil
	}

	var data []zap.Field
	for i := 0; i < len(keyvals); i += 2 {
		data = append(data, zap.Any(fmt.Sprint(keyvals[i]), keyvals[i+1]))
	}

	_logger := l.log
	switch level {
	case log.LevelDebug:
		_logger.Debug("", data...)
	case log.LevelInfo:
		_logger.Info("", data...)
	case log.LevelWarn:
		_logger.Warn("", data...)
	case log.LevelError:
		_logger.Error("", data...)
	case log.LevelFatal:
		_logger.Fatal("", data...)
	}
	return nil
}
