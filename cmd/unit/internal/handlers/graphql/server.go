package graphql

import (
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/devdammit/shekel/pkg/log"
	"github.com/devdammit/shekel/pkg/middleware"
	"github.com/devdammit/shekel/pkg/service"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func NewServer(addr string, shutdownTimeout time.Duration, resolver *Resolver) *service.HTTPServer {
	router := httprouter.New()

	srv := handler.NewDefaultServer(NewExecutableSchema(Config{Resolvers: resolver}))

	router.Handle("GET", "/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		h := playground.AltairHandler("GraphQL playground", "/query")
		h.ServeHTTP(w, r)
	})

	router.Handle("POST", "/query", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		srv.ServeHTTP(w, r)
	})

	return service.NewHTTPServer(
		addr,
		shutdownTimeout,
		alice.New(
			middleware.Logger(log.Default()),
			middleware.Recover(log.Default(), middleware.LogPanicRequest),
			middleware.Tracing(),
		).Then(router),
	)
}
