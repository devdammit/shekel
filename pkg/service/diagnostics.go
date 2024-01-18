package service

import (
	"net/http"
	"net/http/pprof"

	"github.com/devdammit/shekel/pkg/encoding/json"
	"github.com/devdammit/shekel/pkg/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	serviceInfo map[string]string
)

func SetInfo(info map[string]string) {
	serviceInfo = info
}

type DiagnosticServerConfig struct {
	DocsDir             string
	ConfigClientHandler http.HandlerFunc
	StateHandler        http.HandlerFunc
}

func StartDiagnosticsServer(addr string) {
	startDiagnosticsServer(addr, "", nil, nil)
}

func StartDiagnosticsServerWithConfig(addr string, cfg DiagnosticServerConfig) {
	if addr != "" {
		startDiagnosticsServer(addr, cfg.DocsDir, cfg.ConfigClientHandler, cfg.StateHandler)
	}
}

func startDiagnosticsServer(addr, docs string, configsHandler http.HandlerFunc, stateHandler http.HandlerFunc) {
	jsonInfo, err := json.Marshal(serviceInfo)
	if err != nil {
		panic(err)
	}
	mux := http.NewServeMux()

	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/info", func(w http.ResponseWriter, _ *http.Request) {
		if _, err := w.Write(jsonInfo); err != nil {
			log.With(log.Err(err)).Error("failed to write http response")
		}
	})

	// docs
	if docs != "" {
		mux.Handle("/docs/", http.FileServer(http.Dir(docs)))
	}
	if stateHandler != nil {
		mux.Handle("/state/", stateHandler)
	}
	go NewHTTPServer(addr, 0, mux).Start()
}
