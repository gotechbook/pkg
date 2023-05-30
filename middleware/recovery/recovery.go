package recovery

import (
	"context"
	"github.com/gotechbook/pkg/errors"
	"github.com/gotechbook/pkg/logger"
	"github.com/gotechbook/pkg/middleware"
	"runtime"
)

type HandlerFunc func(ctx context.Context, req, err interface{}) error

type Option func(*options)

type options struct {
	handler HandlerFunc
}

func WithHandler(h HandlerFunc) Option {
	return func(o *options) {
		o.handler = h
	}
}

func Recovery(opts ...Option) middleware.Middleware {
	op := options{func(ctx context.Context, req, err interface{}) error {
		return errors.InternalServer("UNKNOWN", "unknown request error")
	}}
	for _, o := range opts {
		o(&op)
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			defer func() {
				if re := recover(); re != nil {
					buf := make([]byte, 64<<10)
					n := runtime.Stack(buf, false)
					buf = buf[:n]
					logger.Errorf("%v: %+v\n%s\n", re, req, buf)
					err = op.handler(ctx, req, re)
				}
			}()
			return handler(ctx, req)
		}
	}
}
