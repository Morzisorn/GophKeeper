// Package logger i used to log info
package logger

import (
	"go.uber.org/zap"
)

var Log *zap.Logger = zap.NewNop()

// Init creates new logger
// Should be used once when app starts
func Init() error {
	cfg := zap.NewProductionConfig()

	cfg.Level.SetLevel(zap.InfoLevel)

	cfg.EncoderConfig.StacktraceKey = ""

	zl, err := cfg.Build()
	if err != nil {
		return err
	}

	Log = zl
	return nil
}
