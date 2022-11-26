package grpc

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"sync"
	"time"

	"github.com/weiqiangxu/common-config/logger"
	"github.com/weiqiangxu/user/net/internal/host"
	"github.com/weiqiangxu/user/net/transport"

	grpcRecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
)

var _ transport.Server = (*Server)(nil)

const HealthcheckService = "grpc.health.v1.Health"

type ServerOption func(o *Server)

func UnaryInterceptor(in ...grpc.UnaryServerInterceptor) ServerOption {
	return func(s *Server) {
		s.interceptor = in
	}
}

// Address with server address.
func Address(addr string) ServerOption {
	return func(s *Server) {
		s.address = addr
	}
}

type Server struct {
	*grpc.Server
	ctx         context.Context
	listener    net.Listener
	once        sync.Once
	err         error
	network     string
	address     string
	endpoint    *url.URL
	timeout     time.Duration
	interceptor []grpc.UnaryServerInterceptor
	grpcOpts    []grpc.ServerOption
	health      *health.Server

	tracing  bool
	recovery bool
}

func NewServer(opts ...ServerOption) *Server {
	srv := &Server{
		network: "tcp",
		address: ":0",
		timeout: 1 * time.Second,
		health:  health.NewServer(),
	}
	for _, o := range opts {
		o(srv)
	}
	var unaryInterceptors []grpc.UnaryServerInterceptor
	var streamInterceptors []grpc.StreamServerInterceptor
	if len(srv.interceptor) > 0 {
		unaryInterceptors = append(unaryInterceptors, srv.interceptor...)
	}
	if srv.recovery {
		grpcRecoveryOpts := []grpcRecovery.Option{
			grpcRecovery.WithRecoveryHandlerContext(func(ctx context.Context, p interface{}) (err error) {
				md, ok := metadata.FromIncomingContext(ctx)
				if ok {
					logger.Errorf("gRPC metadata: %v panic info: %v", md, p)
				} else {
					logger.Error(p)
				}
				return fmt.Errorf("context panic triggered: %v", p)
			}),
		}
		unaryInterceptors = append(unaryInterceptors, grpcRecovery.UnaryServerInterceptor(grpcRecoveryOpts...))
		streamInterceptors = append(streamInterceptors, grpcRecovery.StreamServerInterceptor(grpcRecoveryOpts...))
	}
	if srv.tracing {
		unaryInterceptors = append(unaryInterceptors, otelgrpc.UnaryServerInterceptor())
		streamInterceptors = append(streamInterceptors, otelgrpc.StreamServerInterceptor())
	}
	grpcOpts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(unaryInterceptors...),
		grpc.ChainStreamInterceptor(streamInterceptors...),
	}
	if len(srv.grpcOpts) > 0 {
		grpcOpts = append(grpcOpts, srv.grpcOpts...)
	}
	srv.Server = grpc.NewServer(grpcOpts...)
	srv.health.SetServingStatus(HealthcheckService, grpc_health_v1.HealthCheckResponse_SERVING)
	grpc_health_v1.RegisterHealthServer(srv.Server, srv.health)
	reflection.Register(srv.Server)
	return srv
}

func (s *Server) Start(ctx context.Context) error {
	if _, err := s.Endpoint(); err != nil {
		return err
	}
	s.ctx = ctx
	logger.Infof("[gRPC] server listening on: %s", s.listener.Addr().String())
	s.health.Resume()
	return s.Serve(s.listener)
}

func (s *Server) Stop(ctx context.Context) error {
	s.GracefulStop()
	s.health.Shutdown()
	logger.Info("[gRPC] server stopping")
	return nil
}

// Endpoint return a real address to registry endpoint.
// examples:
//
//	grpc://127.0.0.1:9000?isSecure=false
func (s *Server) Endpoint() (*url.URL, error) {
	s.once.Do(func() {
		lis, err := net.Listen(s.network, s.address)
		if err != nil {
			s.err = err
			return
		}
		addr, err := host.Extract(s.address, s.listener)
		if err != nil {
			lis.Close()
			s.err = err
			return
		}
		s.listener = lis
		s.endpoint = &url.URL{Scheme: "grpc", Host: addr}
	})
	if s.err != nil {
		return nil, s.err
	}
	return s.endpoint, nil
}
