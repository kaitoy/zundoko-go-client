package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var sugaredLogger *zap.SugaredLogger

func init() {
	sugaredLogger = zap.NewNop().Sugar()
}

// Init returns the logger.
func Init(level zapcore.Level) {
	logConfig := zap.Config{
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		Level:            zap.NewAtomicLevelAt(level),
		Encoding:         "console",
		EncoderConfig: zapcore.EncoderConfig{
			LevelKey:         "level",
			TimeKey:          "time",
			MessageKey:       "msg",
			CallerKey:        "caller",
			EncodeTime:       zapcore.ISO8601TimeEncoder,
			EncodeLevel:      zapcore.LowercaseLevelEncoder,
			EncodeCaller:     zapcore.ShortCallerEncoder,
			ConsoleSeparator: " ",
		},
	}

	logger, err := logConfig.Build()
	if err != nil {
		panic(err)
	}

	sugaredLogger = logger.Sugar()
}

// GetLogger returns the logger.
func GetLogger() *zap.SugaredLogger {
	return sugaredLogger
}
