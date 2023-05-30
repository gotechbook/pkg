package middleware

import "context"

type Handler func(ctx context.Context, req interface{}) (interface{}, error)

type Middleware func(Handler) Handler

func Chain(m ...Middleware) Middleware {
	return func(handler Handler) Handler {
		for i := len(m) - 1; i >= 0; i-- {
			handler = m[i](handler)
		}
		return handler
	}
}
