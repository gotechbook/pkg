package grpc

import (
	"context"
	ic "github.com/gotechbook/pkg/context"
	"github.com/gotechbook/pkg/middleware"
	"github.com/gotechbook/pkg/transport"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func (s *Server) streamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx, cancel := ic.Merge(ss.Context(), s.context)
		defer cancel()
		md, _ := metadata.FromIncomingContext(ctx)
		replyHeader := metadata.MD{}
		ctx = transport.NewServerContext(ctx, &Transport{
			endpoint:    s.endpoint.String(),
			operation:   info.FullMethod,
			reqHeader:   headerCarrier(md),
			replyHeader: headerCarrier(replyHeader),
		})
		err := handler(srv, ss)
		if len(replyHeader) > 0 {
			_ = grpc.SetHeader(ctx, replyHeader)
		}
		return err
	}
}

func (s *Server) unaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		ctx, cancel := ic.Merge(ctx, s.context)
		defer cancel()

		md, _ := metadata.FromIncomingContext(ctx)
		replyHeader := metadata.MD{}
		tr := &Transport{
			operation:   info.FullMethod,
			reqHeader:   headerCarrier(md),
			replyHeader: headerCarrier(replyHeader),
		}

		if s.endpoint != nil {
			tr.endpoint = s.endpoint.String()
		}

		ctx = transport.NewServerContext(ctx, tr)
		if s.timeout > 0 {
			ctx, cancel = context.WithTimeout(ctx, s.timeout)
			defer cancel()
		}

		h := func(ctx context.Context, req interface{}) (interface{}, error) {
			return handler(ctx, req)
		}
		if next := s.middleware.Matcher(tr.Operation()); len(next) > 0 {
			h = middleware.Chain(next...)(h)
		}

		reply, err := h(ctx, req)
		if len(replyHeader) > 0 {
			_ = grpc.SetHeader(ctx, replyHeader)
		}
		return reply, err
	}
}

func ExtractIncoming(ctx context.Context) metadata.MD {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return metadata.MD(metadata.Pairs())
	}
	return metadata.MD(md)
}
