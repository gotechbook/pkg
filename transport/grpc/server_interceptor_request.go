package grpc

import (
	"context"
	"github.com/google/uuid"
	"github.com/gotechbook/pkg/logger"
	"google.golang.org/grpc"
	"time"
)

const clientRequestIDKey = "X-Request-ID"

func RequestLogUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		starTime := time.Now()
		resp, err = handler(ctx, req)
		logger.Infow("method", info.FullMethod,
			"id", ExtractIncoming(ctx).Get(clientRequestIDKey),
			"request", req,
			"response", resp,
			"err", err,
			"duration", time.Now().Sub(starTime),
		)
		return
	}
}

func RequestLogStreamInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		lss := &streamLog{
			ServerStream: ss,
			fullMethod:   info.FullMethod,
			cid:          ExtractIncoming(ss.Context()).Get(clientRequestIDKey)[0],
			rid:          uuid.New().String(),
		}
		return handler(srv, lss)
	}
}

type streamLog struct {
	grpc.ServerStream
	fullMethod string
	rid        string
	cid        string
}

func (r *streamLog) SendMsg(m any) error {
	if err := r.ServerStream.SendMsg(m); err != nil {
		return err
	}
	logger.Infow("method", r.fullMethod,
		"id", r.rid,
		"cid", r.cid,
		"request", m,
	)
	return nil
}

func (r *streamLog) RecvMsg(m any) error {
	if err := r.ServerStream.RecvMsg(m); err != nil {
		return err
	}
	logger.Infow("method", r.fullMethod,
		"id", r.rid,
		"cid", r.cid,
		"response", m,
	)
	return nil
}
