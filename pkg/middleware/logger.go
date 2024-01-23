package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/devdammit/shekel/pkg/log"
)

// Logger middleware setup logger and logs all requests.
func Logger(logger log.Logger) func(http.Handler) http.Handler {
	return filterLogger(logger, nil)
}

// LoggerWithFilter middleware setup logger and logs requests that satisfy conditions set in filter.
func LoggerWithFilter(logger log.Logger, filter func(o ResponseObserver, r *http.Request) bool) func(http.Handler) http.Handler {
	return filterLogger(logger, filter)
}

type ResponseObserver interface {
	StatusCode() int
}

func filterLogger(logger log.Logger, filter func(o ResponseObserver, r *http.Request) bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			o := &responseObserver{ResponseWriter: w}
			start := time.Now()
			next.ServeHTTP(o, r)
			if filter == nil || filter(o, r) {
				logger.WithContext(r.Context()).With([]log.Field{
					log.String("request", fmt.Sprintf("%s %s %s", r.Method, r.RequestURI, r.Proto)),
					log.String("request_method", r.Method),
					log.String("request_uri", r.RequestURI),
					log.String("request_proto", r.Proto),
					log.Int("request_duration_ms", int(time.Since(start).Seconds()*1000)),
					log.Int("status", o.status),
					log.Int64("content_length", o.written),
					log.String("real_ip", r.Header.Get("X-Real-IP")),
					log.String("proxy_add_x_forwarded_for", r.Header.Get("X-Forwarded-For")),
					log.String("remote_addr", r.RemoteAddr),
					log.String("http_referrer", r.Referer()),
					log.String("http_user_agent", r.UserAgent()),
				}...).Info("http request handled")
			}
		})
	}
}
