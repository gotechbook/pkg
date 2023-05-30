package grpc

import (
	"crypto/tls"
	"github.com/gotechbook/pkg/middleware"
	"google.golang.org/grpc"
	"time"
)

type ClientOption func(o *Client)

func WithClientEndpoint(endpoint string) ClientOption {
	return func(o *Client) {
		o.endpoint = endpoint
	}
}

func WithClientMiddleware(m ...middleware.Middleware) ClientOption {
	return func(o *Client) {
		o.middleware.Use(m...)
	}
}

func WithClientStreamInterceptor(in ...grpc.StreamClientInterceptor) ClientOption {
	return func(o *Client) {
		o.streamInterceptor = in
	}
}

func WithClientTimeout(timeout time.Duration) ClientOption {
	return func(o *Client) {
		o.timeout = timeout
	}
}

func WithClientTLSConfig(c *tls.Config) ClientOption {
	return func(o *Client) {
		o.tlsConf = c
	}
}

func WithClientUnaryInterceptor(in ...grpc.UnaryClientInterceptor) ClientOption {
	return func(o *Client) {
		o.unaryInterceptor = in
	}
}
