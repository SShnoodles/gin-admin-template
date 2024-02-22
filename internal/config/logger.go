package config

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

var Logger *zap.Logger
var Log *zap.SugaredLogger

func init() {
	lumberjackLogger := &lumberjack.Logger{
		Filename:   AppConfig.Logging.File.Path + "/" + AppConfig.Logging.File.Name,
		MaxSize:    100, // MB
		MaxBackups: 7,
		MaxAge:     7, // days
		Compress:   true,
	}

	file := zapcore.AddSync(lumberjackLogger)
	console := zapcore.AddSync(os.Stdout)
	write := zapcore.NewMultiWriteSyncer(file, console)
	coreLogger := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()),
		write,
		zap.InfoLevel,
	)

	log := zap.New(coreLogger, zap.AddCaller())
	defer log.Sync()

	Logger = log
	Log = log.Sugar()
}
