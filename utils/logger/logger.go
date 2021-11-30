package logger

import (
	"os"

	"github.com/natefinch/lumberjack"
	"missevan-fm/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Init(conf *config.LogConfig) (err error) {
	level := new(zapcore.Level) // set level
	err = level.UnmarshalText([]byte(conf.Level))
	if err != nil {
		return
	}

	writeSyncer := getLogWriter(conf)
	defaultEncoder := getLogEncoder()

	var core zapcore.Core
	if conf.Level == "debug" {
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

func getLogEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.0000")
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

func getLogWriter(conf *config.LogConfig) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   conf.File,
		MaxSize:    conf.MaxSize,
		MaxBackups: conf.MaxBackups,
		MaxAge:     conf.MaxAge,
	}
	return zapcore.AddSync(lumberJackLogger)
}
