package grpc

import (
	"context"
	"crypto/tls"
	"github.com/gotechbook/pkg/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	secure "google.golang.org/grpc/credentials/insecure"
	"time"
)

type Client struct {
	endpoint          string
	timeout           time.Duration
	tlsConf           *tls.Config
	middleware        middleware.Matcher
	grpcClientOpts    []grpc.DialOption
	unaryInterceptor  []grpc.UnaryClientInterceptor
	streamInterceptor []grpc.StreamClientInterceptor
}

func Dial(ctx context.Context, opts ...ClientOption) (*grpc.ClientConn, error) {
	return dial(ctx, false, opts...)
}

func DialInsecure(ctx context.Context, opts ...ClientOption) (*grpc.ClientConn, error) {
	return dial(ctx, true, opts...)
}

func dial(ctx context.Context, insecure bool, opts ...ClientOption) (*grpc.ClientConn, error) {
	client := &Client{
		timeout: 2000 * time.Millisecond,
	}
	for _, o := range opts {
		o(client)
	}

	unaryInterceptor := []grpc.UnaryClientInterceptor{
		client.unaryClientInterceptor(),
	}
	streamInterceptor := []grpc.StreamClientInterceptor{
		client.streamClientInterceptor(),
	}

	if len(client.unaryInterceptor) > 0 {
		unaryInterceptor = append(unaryInterceptor, client.unaryInterceptor...)
	}
	if len(client.streamInterceptor) > 0 {
		streamInterceptor = append(streamInterceptor, client.streamInterceptor...)
	}

	grpcClientOption := []grpc.DialOption{
		grpc.WithChainUnaryInterceptor(unaryInterceptor...),
		grpc.WithChainStreamInterceptor(streamInterceptor...),
	}
	if insecure {
		grpcClientOption = append(grpcClientOption, grpc.WithTransportCredentials(secure.NewCredentials()))
	}
	if client.tlsConf != nil {
		grpcClientOption = append(grpcClientOption, grpc.WithTransportCredentials(credentials.NewTLS(client.tlsConf)))
	}
	if len(client.grpcClientOpts) > 0 {
		grpcClientOption = append(grpcClientOption, client.grpcClientOpts...)
	}

	return grpc.DialContext(ctx, client.endpoint, grpcClientOption...)
}
