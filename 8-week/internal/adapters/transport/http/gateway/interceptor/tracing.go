package interceptor

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type metadataTextMapWriterCarrier metadata.MD

func (c metadataTextMapWriterCarrier) Set(key, val string) {
	md := metadata.MD(c)
	md.Set(key, val)
}

func TracingInterceptor(tracer opentracing.Tracer) grpc.UnaryClientInterceptor {
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
		span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, tracer, method, ext.SpanKindRPCClient)
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
