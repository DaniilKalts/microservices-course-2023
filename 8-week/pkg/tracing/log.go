package tracing

import (
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"go.uber.org/zap"
)

func LogWarn(logger *zap.Logger, span opentracing.Span, msg string, err error, fields ...zap.Field) {
	logger.Warn(msg, append(fields, zap.Error(err))...)
	span.LogKV("event", "warning", "message", err.Error())
}

func LogError(logger *zap.Logger, span opentracing.Span, msg string, err error, fields ...zap.Field) {
	logger.Error(msg, append(fields, zap.Error(err))...)
	ext.Error.Set(span, true)
	span.LogKV("event", "error", "message", err.Error())
}
