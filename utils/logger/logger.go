package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Init(lvl string) (err error) {
	level := new(zapcore.Level) // set level
	if err = level.UnmarshalText([]byte(lvl)); err != nil {
		return
	}

	writeSyncer := logWriter()
	defaultEncoder := logEncoder()

	var core zapcore.Core
	if lvl == "debug" {
		// use both of two zapcore in dev mode
		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		core = zapcore.NewTee(
			zapcore.NewCore(defaultEncoder, writeSyncer, level),
			zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel),
		)
	} else {
		// only use default zapcore, not showing in standard output
		core = zapcore.NewCore(defaultEncoder, writeSyncer, level)
	}
	logger := zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(logger)
	return
}

func logEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.0000")
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func logWriter() zapcore.WriteSyncer {
	file, err := os.OpenFile("missevan.log", os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic(fmt.Errorf("create log file failed: %s", err))
	}
	return zapcore.AddSync(file)
}
