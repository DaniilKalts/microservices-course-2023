package middleware

import (
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type statusCapturingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *statusCapturingResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *statusCapturingResponseWriter) Flush() {
	if f, ok := w.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
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
