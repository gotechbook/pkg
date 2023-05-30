package grpc

import (
	"context"
	"github.com/gotechbook/pkg/transport"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func (c *Client) streamClientInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		ctx = transport.NewClientContext(ctx, &Transport{
			endpoint:  cc.Target(),
			operation: method,
			reqHeader: headerCarrier{},
		})
		return streamer(ctx, desc, cc, method, opts...)
	}
}

func (c *Client) unaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {

		ctx = transport.NewClientContext(ctx, &Transport{
			endpoint:  cc.Target(),
			operation: method,
			reqHeader: headerCarrier{},
		})

		if c.timeout > 0 {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, c.timeout)
			defer cancel()
		}

		h := func(ctx context.Context, req interface{}) (interface{}, error) {
			if tr, ok := transport.FromClientContext(ctx); ok {
				header := tr.RequestHeader()
				keys := header.Keys()
				ks := make([]string, 0, len(keys))
				for _, k := range keys {
					ks = append(ks, k, header.Get(k))
				}
				ctx = metadata.AppendToOutgoingContext(ctx, ks...)
			}
			return reply, invoker(ctx, method, req, reply, cc, opts...)
		}

		_, err := h(ctx, req)
		return err
	}
}
