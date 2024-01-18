package service

import (
	"context"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/devdammit/shekel/pkg/log"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	ready     bool
	readyLock sync.RWMutex
	alive     bool
	aliveLock sync.RWMutex
)

func SetReady(state bool) {
	readyLock.Lock()
	defer readyLock.Unlock()
	ready = state
}

func IsReady() bool {
	readyLock.RLock()
	defer readyLock.RUnlock()
	return ready
}

func SetAlive(state bool) {
	aliveLock.Lock()
	defer aliveLock.Unlock()
	alive = state
}

func IsAlive() bool {
	aliveLock.RLock()
	defer aliveLock.RUnlock()
	return alive
}

func init() {
	SetAlive(true)
}

type connStates struct {
	sync.RWMutex
	conns  map[net.Conn]http.ConnState
	states map[http.ConnState][]net.Conn
}

func filter(conns []net.Conn, conn net.Conn) []net.Conn {
	ret := conns[:0]
	for _, v := range conns {
		if v != conn {
			ret = append(ret, v)
		}
	}
	if len(conns) != len(ret) {
		conns[len(conns)-1] = nil
	}
	return ret
}

func (c *connStates) setState(conn net.Conn, state http.ConnState) {
	c.Lock()
	defer c.Unlock()
	for st, conns := range c.states {
		c.states[st] = filter(conns, conn)
	}
	if state == http.StateClosed {
		delete(c.conns, conn)
		return
	}
	c.conns[conn] = state
	c.states[state] = append(c.states[state], conn)
}

type HTTPServer struct {
	Server          *http.Server
	shutdownTimeout time.Duration
	done            chan struct{}
}

func NewHTTPServer(addr string, shutdownTimeout time.Duration, router http.Handler) *HTTPServer {
	mux := http.NewServeMux()
	mux.Handle("/", router)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if IsAlive() {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
	})
	mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		if IsReady() {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
	})

	serverLabels := make(map[string]string)
	serverLabels["addr"] = addr
	connCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace:   "http_server",
		Name:        "connection_state_changes_total",
		Help:        "Connection state changes count",
		ConstLabels: serverLabels,
	}, []string{"state"})
	prometheus.MustRegister(connCounter)
	connGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace:   "http_server",
		Name:        "connections_total",
		Help:        "Connections count by states",
		ConstLabels: serverLabels,
	}, []string{"state"})
	prometheus.MustRegister(connGauge)

	states := &connStates{
		conns:  make(map[net.Conn]http.ConnState),
		states: make(map[http.ConnState][]net.Conn),
	}
	// nolint gosec
	// FIXME: enable ReadHeaderTimeout after https://github.com/golang/go/issues/54784 is resolved
	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
		// ReadHeaderTimeout: 10 * time.Second,
		ConnState: func(conn net.Conn, state http.ConnState) {
			connCounter.WithLabelValues(state.String()).Inc()
			states.setState(conn, state)
			states.RLock()
			for state, conns := range states.states {
				connGauge.WithLabelValues(state.String()).Set(float64(len(conns)))
			}
			states.RUnlock()
		},
	}

	return &HTTPServer{
		shutdownTimeout: shutdownTimeout,
		Server:          srv,
		done:            make(chan struct{}),
	}
}

func (s *HTTPServer) Start() {
	go func() {
		defer close(s.done)
		log.With(log.String("address", s.Server.Addr)).Info("starting http server on address")
		if err := s.Server.ListenAndServe(); err != http.ErrServerClosed {
			log.With(log.Err(err)).Fatal("http server failure")
		}
		log.With(log.String("address", s.Server.Addr)).Info("http server stopped listening")
	}()
}

func (s *HTTPServer) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	if err := s.Server.Shutdown(ctx); err != nil {
		log.With(log.Err(err)).Error("http shutdown error")
	}
	log.With(log.String("address", s.Server.Addr)).Info("http server stopped")
	<-s.done
	cancel()
}
