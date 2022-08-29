package testutil

import (
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

// NewLogObserver is use in test codes only.
func NewLogObserver(_ *testing.T, level zapcore.Level) (*observer.ObservedLogs, *zap.Logger) {
	c, observedLogs := observer.New(level)
	logger := zap.New(c)
	zap.ReplaceGlobals(logger)

	return observedLogs, logger
}
