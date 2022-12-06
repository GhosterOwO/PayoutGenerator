package util

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var LogLevel = zapcore.InfoLevel

func NewLogger() *zap.Logger {
	lowPriority := zap.LevelEnablerFunc(func(lv zapcore.Level) bool {
		return lv >= LogLevel
	})

	jsonEnc := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	stdCore := zapcore.NewCore(jsonEnc, zapcore.Lock(os.Stdout), lowPriority)

	return zap.New(stdCore).WithOptions(zap.AddCaller())
}

func NewFileLogger(logName string) *zap.Logger {
	pe := zap.NewProductionEncoderConfig()

	fileEncoder := zapcore.NewJSONEncoder(pe)

	pe.EncodeTime = zapcore.ISO8601TimeEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(pe)

	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logName,
		MaxSize:    50, // megabytes
		MaxBackups: 30,
		MaxAge:     28, // days
	})

	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, w, LogLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), LogLevel),
	)
	return zap.New(core)
}

func NewELKLogger(address string) *zap.Logger {
	w := NewRedisWriter(address, "logs")

	lowPriority := zap.LevelEnablerFunc(func(lv zapcore.Level) bool {
		return lv >= zapcore.DebugLevel
	})

	// 使用 JSON 格式日誌
	jsonEnc := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	stdCore := zapcore.NewCore(jsonEnc, zapcore.Lock(os.Stdout), lowPriority)

	// addSync 將 io.Writer 裝飾爲 WriteSyncer
	// 故只需要一個實現 io.Writer 接口的對象即可
	syncer := zapcore.AddSync(w)
	redisCore := zapcore.NewCore(jsonEnc, syncer, lowPriority)

	// 集成多個 core
	core := zapcore.NewTee(stdCore, redisCore)

	// logger 輸出到 console 且標識調用代碼行
	return zap.New(core).WithOptions(zap.AddCaller())
}
