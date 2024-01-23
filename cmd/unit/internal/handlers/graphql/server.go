package graphql

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/devdammit/shekel/pkg/log"
	"github.com/devdammit/shekel/pkg/middleware"
	"github.com/devdammit/shekel/pkg/service"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"net/http"
	"time"
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

//func NewServer(addr string, shutdownTimeout time.Duration, srvs Services) *service.HTTPServer {
//	router := httprouter.New()
//
//	createPositionEndpoint := createposition.NewEndpoint(srvs.PositionService)
//	createMenuEndpoint := createmenu.NewEndpoint(srvs.MenuService)
//	menuTodayEndpoint := menutoday.NewEndpoint(srvs.MenuService)
//	openEndpoint := open.NewEndpoint(srvs.WorkService)
//	closeEndpoint := close.NewEndpoint(srvs.WorkService)
//
//	router.POST("/api/v1/positions", createPositionEndpoint.Handler())
//	router.POST("/api/v1/menus", createMenuEndpoint.Handler())
//	router.GET("/api/v1/menus/today", menuTodayEndpoint.Handler())
//	router.POST("/api/v1/kitchen/open", openEndpoint.Handler())
//	router.POST("/api/v1/kitchen/close", closeEndpoint.Handler())
//
//	return service.NewHTTPServer(
//		addr,
//		shutdownTimeout,
//		alice.New(
//			middleware.Logger(log.Default()),
//			middleware.Recover(log.Default(), middleware.LogPanicRequest),
//			middleware.Tracing(),
//		).Then(router),
//	)
//}

//func (s Server) Start() {
//	log.Info("starting graphql server")
//
//	srv := handler.NewDefaultServer(NewExecutableSchema(Config{Resolvers: s.resolver}))
//
//	mux := http.NewServeMux()
//	mux.Handle("/", playground.AltairHandler("GraphQL playground", "/query"))
//	mux.Handle("/query", srv)
//
//	serv := &http.Server{
//		Addr:    ":" + s.port,
//		Handler: mux,
//	}
//
//	s.srv = serv
//
//	log.Info(fmt.Sprintf("connect to http://localhost:%s/ for GraphQL playground", s.port))
//
//	go func() {
//		if err := serv.ListenAndServe(); err != nil {
//			log.With(log.Err(err)).Fatal("graphql server failure")
//		}
//	}()
//}
//
//func (s Server) Stop() {
//	if err := s.srv.Shutdown(context.Background()); err != nil {
//		log.With(log.Err(err)).Fatal("graphql server shutdown failure")
//	}
//
//	log.Info("graphql server stopped")
//}
