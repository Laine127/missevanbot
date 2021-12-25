package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Init(path, lvl string) (err error) {
	level := new(zapcore.Level) // set level
	if err = level.UnmarshalText([]byte(lvl)); err != nil {
		return
	}

	writeSyncer := logWriteSyncer(path)
	encoder := logEncoder()

	var core zapcore.Core
	if lvl == "debug" {
		// use both of two zap core in dev mode
		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		core = zapcore.NewTee(
			zapcore.NewCore(encoder, writeSyncer, level),
			zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel),
		)
	} else {
		// only use default zap core, not showing in standard output
		core = zapcore.NewCore(encoder, writeSyncer, level)
	}
	logger := zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(logger)
	return
}

func logEncoder() zapcore.Encoder {
	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.0000")
	cfg.EncodeLevel = zapcore.CapitalLevelEncoder
	cfg.EncodeDuration = zapcore.SecondsDurationEncoder
	cfg.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewConsoleEncoder(cfg)
}

func logWriteSyncer(path string) zapcore.WriteSyncer {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic(fmt.Errorf("create log file failed: %s", err))
	}
	return zapcore.AddSync(file)
}
