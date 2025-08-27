package logger

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitializeDevelopmentLogger(serviceName string) (*ContextLogger, error) {
	logConfig := NewLogConfig()
	logConfig.InitializeLogConfig(serviceName)
	return NewLogger(*logConfig)
}

func InitializeProductionLogger(serviceName string) (*ContextLogger, error) {
	logConfig := NewLogConfig()
	logConfig.InitializeLogConfig(serviceName)
	return NewLogger(*logConfig)
}

func NewLogger(config LogConfig) (*ContextLogger, error) {
	var zapConfig zap.Config

	switch config.Environment {
	case "production":
		zapConfig = zap.NewProductionConfig()
		zapConfig.EncoderConfig.TimeKey = "timestamp"
		zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	default:
		zapConfig = zap.NewDevelopmentConfig()
		zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	level, err := zap.ParseAtomicLevel(config.Level)
	if err != nil {
		return nil, fmt.Errorf("invalid log level %s: %w", config.Level, err)
	}
	zapConfig.Level = level

	zapConfig.InitialFields = map[string]interface{}{
		"service":     config.ServiceName,
		"version":     config.Version,
		"environment": config.Environment,
	}

	logger, err := zapConfig.Build(zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	if err != nil {
		return nil, fmt.Errorf("failed to build logger: %w", err)
	}

	return &ContextLogger{Logger: logger}, nil
}

// WithContext adds context-based fields to the logger
func (l *ContextLogger) WithContext(ctx context.Context) *zap.Logger {
	fields := []zap.Field{}

	if requestID, ok := ctx.Value(requestIDKey).(string); ok && requestID != "" {
		fields = append(fields, zap.String("request_id", requestID))
	}

	if traceID, ok := ctx.Value(traceIDKey).(string); ok && traceID != "" {
		fields = append(fields, zap.String("trace_id", traceID))
	}

	if userID, ok := ctx.Value(userIDKey).(string); ok && userID != "" {
		fields = append(fields, zap.String("user_id", userID))
	}

	if len(fields) == 0 {
		return l.Logger
	}

	return l.Logger.With(fields...)
}

// LogFunction is a decorator that logs function entry and exit
func (l *ContextLogger) LogFunction(ctx context.Context, funcName string) func() {
	logger := l.WithContext(ctx)
	start := time.Now()

	logger.Debug("Function started", zap.String("function", funcName))

	return func() {
		duration := time.Since(start)
		logger.Info("Function completed",
			zap.String("function", funcName),
			zap.Duration("duration", duration),
		)
	}
}

// LogError logs an error with optional context
func (l *ContextLogger) LogError(ctx context.Context, err error, msg string, fields ...zap.Field) {
	logger := l.WithContext(ctx)
	allFields := append(fields, zap.Error(err))
	logger.Error(msg, allFields...)
}

// LogHTTPRequest logs HTTP request details
func (l *ContextLogger) LogHTTPRequest(ctx context.Context, method, path string, statusCode int, duration time.Duration, fields ...zap.Field) {
	logger := l.WithContext(ctx)
	allFields := append(fields,
		zap.String("method", method),
		zap.String("path", path),
		zap.Int("status_code", statusCode),
		zap.Duration("duration", duration),
	)
	logger.Info("HTTP request completed", allFields...)
}

// Sync flushes any buffered log entries - call this before application exit
func (l *ContextLogger) Sync() error {
	return l.Logger.Sync()
}
