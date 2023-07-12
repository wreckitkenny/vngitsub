package utils

import (
	"os"
	"strconv"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ConfigZap configures default logger configuration
func ConfigZap() *zap.SugaredLogger{
    level, _ := strconv.Atoi(os.Getenv("LOGLEVEL"))

	cfg := zap.Config{
        Encoding:    "console",
        Level:       zap.NewAtomicLevelAt(zapcore.Level(level)),
        OutputPaths: []string{"stderr"},

        EncoderConfig: zapcore.EncoderConfig{
            MessageKey: "msg",
			TimeKey: "time",
            LevelKey: "level",
			CallerKey: "caller",
        	EncodeCaller: zapcore.ShortCallerEncoder,
			EncodeLevel: CustomLevelEncoder,
			EncodeTime: SyslogTimeEncoder,
        },
    }

    logger, _ := cfg.Build()
	return logger.Sugar()
}

// SyslogTimeEncoder - Time encoder
func SyslogTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
    enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

// CustomLevelEncoder - Level encoder
func CustomLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
    enc.AppendString("[" + level.CapitalString() + "]")
}
