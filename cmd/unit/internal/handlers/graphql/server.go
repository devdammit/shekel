package graphql

import (
	"fmt"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/devdammit/shekel/pkg/log"
	"net/http"
)

type Server struct {
	port string

	done chan struct{}
}

func NewServer(port string) *Server {
	return &Server{
		port: port,
		done: make(chan struct{}),
	}
}

func (s Server) Start() {
	log.Info("starting graphql server")

	srv := handler.NewDefaultServer(NewExecutableSchema(Config{Resolvers: &Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Info(fmt.Sprintf("connect to http://localhost:%s/ for GraphQL playground", s.port))

	go func() {
		defer close(s.done)

		if err := http.ListenAndServe(":"+s.port, nil); err != nil {
			log.With(log.Err(err)).Fatal("graphql server failure")
		}
	}()
}

func (s Server) Stop() {
	<-s.done

	log.Info("graphql server stopped")
}
