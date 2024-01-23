package middleware

import (
	"errors"
	"net/http"

	"github.com/devdammit/shekel/pkg/log"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

func Tracing() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			spanCtx, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
			if err != nil && !errors.Is(err, opentracing.ErrSpanContextNotFound) {
				log.With(log.Err(err), log.String("route", r.URL.Path)).Debug("failed to extract tracing context")
			}
			span := opentracing.StartSpan(r.URL.Path, ext.RPCServerOption(spanCtx))
			defer span.Finish()

			ctx := opentracing.ContextWithSpan(r.Context(), span)
			err = opentracing.GlobalTracer().Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(w.Header()))
			if err != nil {
				log.With(log.Err(err), log.String("route", r.URL.Path)).Debug("failed to inject tracing info")
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
