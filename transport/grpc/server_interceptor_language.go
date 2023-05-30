package grpc

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const languageKey = "X-Device-Accept-Language"

func LanguageUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		val := ExtractIncoming(ctx).Get(languageKey)
		if val == nil || val[0] == "" {
			return handler(ctx, req)
		}
		ctx = context.WithValue(ctx, "i18n", ExtractIncoming(ctx).Get(languageKey)[0])
		ctx = metadata.AppendToOutgoingContext(ctx, "i18n", ExtractIncoming(ctx).Get(languageKey)[0])
		return handler(ctx, req)
	}
}

func LanguageStreamInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		val := ExtractIncoming(ss.Context()).Get(languageKey)
		if val == nil || val[0] == "" {
			return handler(srv, ss)
		}
		ptx := ss.Context()
		ntx := context.WithValue(ptx, "i18n", ExtractIncoming(ss.Context()).Get(languageKey)[0])
		ntx = metadata.AppendToOutgoingContext(ntx, "i18n", ExtractIncoming(ss.Context()).Get(languageKey)[0])
		newStream := WrapServerStream(ss)
		newStream.WrappedContext = ntx
		return handler(srv, newStream)
	}
}
