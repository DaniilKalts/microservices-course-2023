package gateway

import (
	"context"
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type statusCapturingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *statusCapturingResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func WithTracing(next http.Handler, tracer opentracing.Tracer) http.Handler {
	if tracer == nil {
		tracer = opentracing.GlobalTracer()
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		operationName := r.Method + " " + r.URL.Path

		spanContext, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
		opts := []opentracing.StartSpanOption{ext.SpanKindRPCServer, opentracing.Tag{Key: "component", Value: "http"}}
		if spanContext != nil {
			opts = append(opts, ext.RPCServerOption(spanContext))
		}

		span := tracer.StartSpan(operationName, opts...)
		defer span.Finish()

		tw := &statusCapturingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(tw, r.WithContext(opentracing.ContextWithSpan(r.Context(), span)))

		span.SetTag("http.method", r.Method)
		span.SetTag("http.url", r.URL.String())
		span.SetTag("http.status_code", tw.statusCode)
		if tw.statusCode >= http.StatusInternalServerError {
			ext.Error.Set(span, true)
		}
	})
}

type metadataTextMapWriterCarrier metadata.MD

func (c metadataTextMapWriterCarrier) Set(key, val string) {
	md := metadata.MD(c)
	md.Set(key, val)
}

func TracingClientInterceptor(tracer opentracing.Tracer) grpc.UnaryClientInterceptor {
	if tracer == nil {
		tracer = opentracing.GlobalTracer()
	}

	return func(
		ctx context.Context,
		method string,
		req, reply any,
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		span, ctx := opentracing.StartSpanFromContext(ctx, method, ext.SpanKindRPCClient)
		defer span.Finish()

		span.SetTag("component", "grpc-client")

		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		} else {
			md = md.Copy()
		}

		if err := tracer.Inject(span.Context(), opentracing.TextMap, metadataTextMapWriterCarrier(md)); err == nil {
			ctx = metadata.NewOutgoingContext(ctx, md)
		}

		err := invoker(ctx, method, req, reply, cc, opts...)
		if err != nil {
			ext.Error.Set(span, true)
			span.LogKV("event", "error", "message", err.Error())
		}

		return err
	}
}
