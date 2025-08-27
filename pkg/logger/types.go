package logger

import (
	"go.uber.org/zap"
)

type contextKey string

type ContextLogger struct {
	*zap.Logger
}

type LogConfig struct {
	Level       string
	Environment string
	ServiceName string
	Version     string
}

const (
	requestIDKey contextKey = "request_id"
	traceIDKey   contextKey = "trace_id"
	userIDKey    contextKey = "user_id"
)
