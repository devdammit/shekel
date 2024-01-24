package fileserve

import (
	"context"
	"errors"
	port "github.com/devdammit/shekel/cmd/unit/internal/ports/repositories"
	"github.com/devdammit/shekel/pkg/log"
	"github.com/devdammit/shekel/pkg/middleware"
	"github.com/devdammit/shekel/pkg/service"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"net/http"
	"time"
)

type Resolver interface {
	GetImage(ctx context.Context, contactID uint64, bankName string) ([]byte, error)
}

func NewServer(addr string, shutdownTimeout time.Duration, resolver Resolver) *service.HTTPServer {
	router := httprouter.New()

	router.Handle(
		"GET",
		"/:contact/:bank",
		func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
			req, err := Parse(context.Background(), request, params)
			if err != nil {
				log.With(log.Err(err)).Error("failed to parse request")
				writer.WriteHeader(http.StatusBadRequest)
				return
			}

			img, err := resolver.GetImage(context.Background(), req.ContactID, req.BankName)
			if err != nil {
				if errors.Is(err, port.ErrNotFound) {
					writer.WriteHeader(http.StatusNotFound)
					return
				}

				log.With(log.Err(err)).Error("failed to get image")
				writer.WriteHeader(http.StatusInternalServerError)
				return
			}

			writer.Header().Set("Content-Type", "image/jpeg")
			writer.WriteHeader(http.StatusOK)
			_, _ = writer.Write(img)
		},
	)

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
