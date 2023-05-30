package grpc

import (
	"crypto/tls"
	"github.com/gotechbook/pkg/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"net"
	"time"
)

type ServerOption func(o *Server)

func WithServerNetWork(network string) ServerOption {
	return func(o *Server) {
		o.network = network
	}
}

func WithServerAddress(address string) ServerOption {
	return func(o *Server) {
		o.address = address
	}
}

func WithServerTimeout(timeout time.Duration) ServerOption {
	return func(o *Server) {
		o.timeout = timeout
	}
}

func WithServerMiddleware(m ...middleware.Middleware) ServerOption {
	return func(o *Server) {
		o.middleware.Use(m...)
	}
}

func WithServerTLSConfig(c *tls.Config) ServerOption {
	return func(o *Server) {
		o.tlsConf = c
	}
}

func WithServerListener(lis net.Listener) ServerOption {
	return func(o *Server) {
		o.listener = lis
	}
}

func WithServerUnaryInterceptor(in ...grpc.UnaryServerInterceptor) ServerOption {
	return func(o *Server) {
		o.unaryInterceptor = in
	}
}

func WithServerStreamInterceptor(in ...grpc.StreamServerInterceptor) ServerOption {
	return func(o *Server) {
		o.streamInterceptor = in
	}
}

func WithServerOption(opts ...grpc.ServerOption) ServerOption {
	return func(o *Server) {
		o.grpcServerOption = opts
	}
}

func WithKeepalive(kep keepalive.EnforcementPolicy, kasp keepalive.ServerParameters) ServerOption {
	return func(o *Server) {
		o.keepaliveEnforcementPolicy = grpc.KeepaliveEnforcementPolicy(kep)
		o.keepaliveParams = grpc.KeepaliveParams(kasp)
	}
}
