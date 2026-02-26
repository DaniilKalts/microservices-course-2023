package middleware

import (
	"errors"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
	bytes      int
}

func Logging(logger *zap.Logger) func(next http.Handler) http.Handler {
	if logger == nil {
		logger = zap.NewNop()
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			recorder := &responseRecorder{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			startedAt := time.Now()
			next.ServeHTTP(recorder, r)
			duration := time.Since(startedAt)

			var requestErr error
			if recorder.statusCode >= http.StatusInternalServerError {
				requestErr = errors.New(http.StatusText(recorder.statusCode))
			}

			errField := zap.NamedError("error", nil)
			if requestErr != nil {
				errField = zap.Error(requestErr)
			}

			fields := []zap.Field{
				zap.String("protocol", "http"),
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status_code", recorder.statusCode),
				zap.String("remote_addr", r.RemoteAddr),
				zap.String("user_agent", r.UserAgent()),
				errField,
				zap.Float64("duration_ms", float64(duration)/float64(time.Millisecond)),
			}

			switch {
			case recorder.statusCode >= http.StatusInternalServerError:
				logger.Error("request completed", fields...)
			case recorder.statusCode >= http.StatusBadRequest:
				logger.Warn("request completed", fields...)
			default:
				logger.Info("request completed", fields...)
			}
		})
	}
}

func (rw *responseRecorder) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *responseRecorder) Write(data []byte) (int, error) {
	written, err := rw.ResponseWriter.Write(data)
	rw.bytes += written
	return written, err
}
