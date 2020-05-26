package log

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewZapDevelopLogger creates a new develop *zap.Logger
func NewZapDevelopLogger() *zap.Logger {
	logger, err := NewDevelopConfig().Build()
	if err != nil {
		log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile | log.LUTC)
		log.Fatalf("could not setup zap logger %v", err)
	}
	// set global logger to develop logger
	_ = zap.ReplaceGlobals(logger)
	return logger
}

// NewDevelopConfig returns the develop logging config
func NewDevelopConfig() zap.Config {
	return zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.DebugLevel),
		Development: true,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding: "console",
		EncoderConfig: zapcore.EncoderConfig{
			// Keys can be anything except the empty string.
			TimeKey:        "T",
			LevelKey:       "L",
			NameKey:        "N",
			CallerKey:      "C",
			MessageKey:     "M",
			StacktraceKey:  "S",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
}
