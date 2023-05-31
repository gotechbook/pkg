package grpc

import (
	"context"
	"github.com/gotechbook/pkg/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strings"
)

const headerAuthorize = "authorization"
const expectedScheme = "bearer"

func AuthorizationUnaryInterceptor(puk string, skipURLs ...string) grpc.UnaryServerInterceptor {
	skipURLsCache := makeSkipURLsCache(skipURLs)
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		if _, ok := skipURLsCache[info.FullMethod]; ok {
			return handler(ctx, req)
		}

		val := ExtractIncoming(ctx).Get(headerAuthorize)
		if val[0] == "" {
			return "", status.Errorf(codes.Unauthenticated, "Request unauthenticated with bearer")
		}
		splits := strings.SplitN(val[0], " ", 2)
		if len(splits) < 2 {
			return "", status.Errorf(codes.Unauthenticated, "Bad authorization string")
		}
		if !strings.EqualFold(splits[0], expectedScheme) {
			return "", status.Errorf(codes.Unauthenticated, "Request unauthenticated with bearer")
		}

		token, err := oauth2.ParseToken(splits[1], puk)
		if err != nil {
			return nil, err
		}
		ctx = context.WithValue(ctx, oauth2.CtxUserTokenKey, token)
		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "bearer "+splits[1])

		return handler(ctx, req)
	}
}
func AuthorizationStreamInterceptor(puk string) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		val := ExtractIncoming(ss.Context()).Get(headerAuthorize)
		if val[0] == "" {
			return status.Errorf(codes.Unauthenticated, "Request unauthenticated with bearer")
		}
		splits := strings.SplitN(val[0], " ", 2)
		if len(splits) < 2 {
			return status.Errorf(codes.Unauthenticated, "Bad authorization string")
		}
		if !strings.EqualFold(splits[0], expectedScheme) {
			return status.Errorf(codes.Unauthenticated, "Request unauthenticated with bearer")
		}
		token, err := oauth2.ParseToken(splits[1], puk)
		if err != nil {
			return err
		}
		ptx := ss.Context()
		ntx := context.WithValue(ptx, oauth2.CtxUserTokenKey, token)
		ntx = metadata.AppendToOutgoingContext(ntx, "authorization", "bearer "+splits[1])

		newStream := WrapServerStream(ss)
		newStream.WrappedContext = ntx

		return handler(srv, newStream)
	}
}

func makeSkipURLsCache(skipURLs []string) map[string]struct{} {
	var cache map[string]struct{}
	if len(skipURLs) > 0 {
		cache = make(map[string]struct{}, len(skipURLs))
		for _, url := range skipURLs {
			cache[url] = struct{}{}
		}
	}
	return cache
}
