package middleware

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type metadataTextMapCarrier metadata.MD

func (c metadataTextMapCarrier) ForeachKey(handler func(key, val string) error) error {
	for key, values := range c {
		for _, value := range values {
			if err := handler(key, value); err != nil {
				return err
			}
		}
	}

	return nil
}

func TracingInterceptor(tracer opentracing.Tracer) grpc.UnaryServerInterceptor {
	if tracer == nil {
		tracer = opentracing.GlobalTracer()
	}

	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		var spanContext opentracing.SpanContext

		if md, ok := metadata.FromIncomingContext(ctx); ok {
			extractedSpanContext, err := tracer.Extract(opentracing.TextMap, metadataTextMapCarrier(md))
			if err == nil {
				spanContext = extractedSpanContext
			}
		}

		opts := []opentracing.StartSpanOption{
			ext.SpanKindRPCServer,
			opentracing.Tag{Key: "component", Value: "grpc"},
			opentracing.Tag{Key: "grpc.method", Value: info.FullMethod},
		}

		if spanContext != nil {
			opts = append(opts, ext.RPCServerOption(spanContext))
		}

		span := tracer.StartSpan(info.FullMethod, opts...)
		defer span.Finish()

		ctx = opentracing.ContextWithSpan(ctx, span)

		resp, err := handler(ctx, req)
		if err != nil {
			ext.Error.Set(span, true)
			span.LogKV("event", "error", "message", err.Error())
		}

		return resp, err
	}
}
