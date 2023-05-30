package grpc

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/gotechbook/pkg/endpoint"
	"github.com/gotechbook/pkg/logger"
	"github.com/gotechbook/pkg/middleware"
	"github.com/gotechbook/pkg/transport"
	"github.com/gotechbook/pkg/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/admin"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"net"
	"net/url"
	"time"
)

var _ transport.Server = (*Server)(nil)
var _ transport.Endpoint = (*Server)(nil)

type Server struct {
	*grpc.Server
	err                        error
	address                    string
	network                    string
	adminClean                 func()
	endpoint                   *url.URL
	listener                   net.Listener
	timeout                    time.Duration
	context                    context.Context
	tlsConf                    *tls.Config
	health                     *health.Server
	middleware                 middleware.Matcher
	grpcServerOption           []grpc.ServerOption
	unaryInterceptor           []grpc.UnaryServerInterceptor
	streamInterceptor          []grpc.StreamServerInterceptor
	keepaliveEnforcementPolicy grpc.ServerOption
	keepaliveParams            grpc.ServerOption
}

func NewServer(opts ...ServerOption) *Server {
	srv := &Server{
		address:    ":0",
		network:    "tcp",
		timeout:    1 * time.Second,
		context:    context.Background(),
		health:     health.NewServer(),
		middleware: middleware.NewMatcher(),
	}
	for _, o := range opts {
		o(srv)
	}

	unaryInterceptor := []grpc.UnaryServerInterceptor{
		srv.unaryServerInterceptor(),
	}
	streamInterceptor := []grpc.StreamServerInterceptor{
		srv.streamServerInterceptor(),
	}

	if len(srv.unaryInterceptor) > 0 {
		unaryInterceptor = append(unaryInterceptor, srv.unaryInterceptor...)
	}

	if len(srv.streamInterceptor) > 0 {
		streamInterceptor = append(streamInterceptor, srv.streamInterceptor...)
	}

	grpcServerOption := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(unaryInterceptor...),
		grpc.ChainStreamInterceptor(streamInterceptor...),
	}
	if srv.tlsConf != nil {
		grpcServerOption = append(grpcServerOption, grpc.Creds(credentials.NewTLS(srv.tlsConf)))
	}
	if len(srv.grpcServerOption) > 0 {
		grpcServerOption = append(grpcServerOption, srv.grpcServerOption...)
	}

	srv.Server = grpc.NewServer(grpcServerOption...)
	grpc_health_v1.RegisterHealthServer(srv.Server, srv.health)
	reflection.Register(srv.Server)

	srv.adminClean, _ = admin.Register(srv.Server)
	return srv
}

func (s *Server) Endpoint() (*url.URL, error) {
	if err := s.listenAndEndpoint(); err != nil {
		return nil, s.err
	}
	return s.endpoint, nil
}

func (s *Server) Start(ctx context.Context) error {
	if err := s.listenAndEndpoint(); err != nil {
		return s.err
	}
	s.context = ctx
	logger.Infow("GRPC", fmt.Sprintf("listening on %s", s.listener.Addr().String()))
	s.health.Resume()
	return s.Serve(s.listener)
}

func (s *Server) Stop(ctx context.Context) error {
	if s.adminClean != nil {
		s.adminClean()
	}
	s.health.Shutdown()
	s.GracefulStop()
	logger.Info("[gRPC] server stopping")
	return nil
}

func (s *Server) listenAndEndpoint() error {
	if s.listener == nil {
		listen, err := net.Listen(s.network, s.address)
		if err != nil {
			s.err = err
			return err
		}
		s.listener = listen
	}

	if s.endpoint == nil {
		addr, err := utils.Extract(s.address, s.listener)
		if err != nil {
			s.err = err
			return err
		}
		s.endpoint = endpoint.NewEndpoint(endpoint.Scheme("grpc", s.tlsConf != nil), addr)
	}
	return s.err
}
