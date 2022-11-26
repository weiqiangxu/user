package http

import (
	"context"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"code.skyhorn.net/backend/infra/logger"

	"github.com/gin-gonic/gin"

	"code.skyhorn.net/backend/infra/gms/config"
	"code.skyhorn.net/backend/infra/gms/transport"
	gin_pprof "github.com/gin-contrib/pprof"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

var _ transport.Server = (*Server)(nil)

type ServerOption func(*Server)

type Server struct {
	gin        *gin.Engine
	httpServer *http.Server
	// listener   net.Listener
	// endpoint   *url.URL
	address string
	network string
	// once       sync.Once
	ms         []gin.HandlerFunc
	prometheus bool
	profile    bool
	tracing    bool
	verbose    bool
	// err        error
}

func WithMiddleware(m ...gin.HandlerFunc) ServerOption {
	return func(server *Server) {
		server.ms = m
	}
}

func WithAddress(addr string) ServerOption {
	return func(server *Server) {
		server.address = addr
	}
}

func WithPrometheus(enablePrometheus bool) ServerOption {
	return func(server *Server) {
		server.prometheus = enablePrometheus
	}
}

func WithProfile(profile bool) ServerOption {
	return func(server *Server) {
		server.profile = profile
	}
}

func WithTracing(tracing bool) ServerOption {
	return func(server *Server) {
		server.tracing = tracing
	}
}

func WithVerbose(verbose bool) ServerOption {
	return func(server *Server) {
		server.verbose = verbose
	}
}

func NewServer(opts ...ServerOption) *Server {
	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()
	g := gin.New()
	srv := &Server{
		network: "tcp",
		address: ":0",
	}
	for _, o := range opts {
		o(srv)
	}
	if srv.prometheus {
		g.GET("metrics", gin.WrapH(promhttp.Handler()))
		prometheus.Unregister(collectors.NewGoCollector())
		prometheus.Unregister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts(prometheus.ProcessCollectorOpts{})))
	}
	if srv.profile {
		gin_pprof.Register(g)
	}
	if srv.tracing {
		g.Use(otelgin.Middleware(config.GetBuildCfg().ServiceName))
		// e.Use(otelecho.Middleware(config.GetBuildCfg().ServiceName, otelecho.WithSkipper(UrlSkipper)))
	}
	if len(srv.ms) > 0 {
		g.Use(srv.ms...)
	}
	g.Use(GinzapWithConfig(&GinLoggerConfig{
		TimeFormat: time.RFC3339,
		UTC:        false,
		SkipPaths:  []string{"/healthz"},
	}))
	g.Use(RecoveryWithZap(true))
	g.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, http.StatusText(http.StatusOK))
	})
	srv.gin = g
	return srv
}

/*func (s *Server) Endpoint() (*url.URL, error) {
	s.once.Do(func() {
		var err error
		s.listener, err = net.Listen(s.network, s.address)
		if err != nil {
			s.err = err
			return
		}
		addr, err := host.Extract(s.address, s.listener)
		if err != nil {
			s.listener.Close()
			s.err = err
			return
		}
		s.endpoint = &url.URL{Scheme: "http", Host: addr}
	})
	if s.err != nil {
		return nil, s.err
	}
	return s.endpoint, nil
}*/

/*func UrlSkipper(c *gin.Context) bool {
	if strings.HasPrefix(c.FullPath(), "/metrics") {
		return true
	}
	if strings.HasPrefix(c.Path(), "/pprof") {
		return true
	}
	if strings.HasPrefix(c.Path(), "/healthz") {
		return true
	}
	return false
}*/

func (s *Server) Server() *gin.Engine {
	return s.gin
}

func (s *Server) Start(ctx context.Context) error {
	/*if _, err := s.Endpoint(); err != nil {
		return err
	}*/
	srv := &http.Server{
		Addr:    s.address,
		Handler: s.gin,
	}
	s.httpServer = srv
	logger.Infof("[HTTP] server listening on: %s", s.address)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	logger.Info("[HTTP] server stopping")
	return s.httpServer.Shutdown(ctx)
}
