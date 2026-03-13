package interceptor

import (
	"context"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type validatable interface {
	Validate() error
}

type multiValidatable interface {
	ValidateAll() error
}

// Interfaces matching protoc-gen-validate error types without importing generated code.
type (
	fieldError interface {
		Field() string
		Reason() string
	}
	multiError interface {
		AllErrors() []error
	}
)

func ValidationInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		_ *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		var validationErr error

		// Prefer ValidateAll to collect every field violation.
		if mv, ok := req.(multiValidatable); ok {
			validationErr = mv.ValidateAll()
		} else if v, ok := req.(validatable); ok {
			validationErr = v.Validate()
		}

		if validationErr != nil {
			st := status.New(codes.InvalidArgument, "validation failed")
			st, _ = st.WithDetails(&errdetails.BadRequest{
				FieldViolations: extractFieldViolations(validationErr),
			})
			return nil, st.Err()
		}

		return handler(ctx, req)
	}
}

func extractFieldViolations(err error) []*errdetails.BadRequest_FieldViolation {
	var errs []error
	if me, ok := err.(multiError); ok {
		errs = me.AllErrors()
	} else {
		errs = []error{err}
	}

	violations := make([]*errdetails.BadRequest_FieldViolation, 0, len(errs))
	for _, e := range errs {
		if fe, ok := e.(fieldError); ok {
			violations = append(violations, &errdetails.BadRequest_FieldViolation{
				Field:       fe.Field(),
				Description: fe.Reason(),
			})
		} else {
			violations = append(violations, &errdetails.BadRequest_FieldViolation{
				Field:       "request",
				Description: e.Error(),
			})
		}
	}

	return violations
}
