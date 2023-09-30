package observ

import (
	"context"
	"net/http"
	"time"
)

var (
	ErrorCtxKey = struct{}{}
)

func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now().UTC()

		fields := []any{
			"method", r.Method,
			"path", r.URL.Path,
			"query", r.URL.Query(),
			"user_agent", r.UserAgent(),
		}

		srw := &statusResponseWriter{ResponseWriter: w}

		defer func() {
			var err error
			if exc := recover(); exc != nil {
				err = exc.(error)
			}

			duration := time.Now().UTC().Sub(start)

			fields = append(fields, "duration", duration)
			fields = append(fields, "status", srw.status)

			if ctxErr := getErrorFromContext(r.Context()); ctxErr != nil {
				err = ctxErr
			}

			if err != nil {
				fields = append(fields, "error", err)
			}

			if stErr, ok := err.(StackTracer); ok {
				fields = append(fields, "stacktrace", stErr.Stacktrace())
			}

			if srw.status >= 200 && srw.status < 500 {
				Log(r.Context(), "HTTP", fields...)
			} else {
				LogError(r.Context(), "HTTP", fields...)
			}
		}()

		next.ServeHTTP(srw, r)
	})
}

func getErrorFromContext(ctx context.Context) error {
	err, _ := ctx.Value(ErrorCtxKey).(error)

	return err
}

type statusResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusResponseWriter) WriteHeader(statusCode int) {
	w.status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
