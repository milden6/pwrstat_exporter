package server

import (
	"context"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	readHeaderTimeout = 500 * time.Millisecond
	readTimeout       = 500 * time.Millisecond
	handlerTimeout    = 1 * time.Second
)

type Server struct {
	httpServer *http.Server
}

func New(addr string) *Server {
	r := setupRouter()

	httpServer := &http.Server{
		Addr:              addr,
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		Handler:           http.TimeoutHandler(r, handlerTimeout, ""),
	}

	return &Server{
		httpServer: httpServer,
	}
}

func (s *Server) Start() error {
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func setupRouter() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("GET /metrics", promhttp.Handler())
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
			<html>
				<head>
					<title>CyberPower UPS Exporter</title>
				</head>
				<body>
					<h1>CyberPower UPS Exporter</h1>
					<p>
						<a href="/metrics">Metrics</a>
					</p>
				</body>
			</html>
		`))
	})

	return mux
}
