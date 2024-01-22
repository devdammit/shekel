package graphql

import (
	"context"
	"fmt"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/devdammit/shekel/pkg/log"
	"net/http"
)

type Server struct {
	port string

	srv *http.Server

	resolver *Resolver
}

func NewServer(port string, resolver *Resolver) *Server {
	return &Server{
		port:     port,
		srv:      &http.Server{},
		resolver: resolver,
	}
}

func (s Server) Start() {
	log.Info("starting graphql server")

	srv := handler.NewDefaultServer(NewExecutableSchema(Config{Resolvers: s.resolver}))

	mux := http.NewServeMux()
	mux.Handle("/", playground.AltairHandler("GraphQL playground", "/query"))
	mux.Handle("/query", srv)

	serv := &http.Server{
		Addr:    ":" + s.port,
		Handler: mux,
	}

	s.srv = serv

	log.Info(fmt.Sprintf("connect to http://localhost:%s/ for GraphQL playground", s.port))

	go func() {
		if err := serv.ListenAndServe(); err != nil {
			log.With(log.Err(err)).Fatal("graphql server failure")
		}
	}()
}

func (s Server) Stop() {
	if err := s.srv.Shutdown(context.Background()); err != nil {
		log.With(log.Err(err)).Fatal("graphql server shutdown failure")
	}

	log.Info("graphql server stopped")
}
