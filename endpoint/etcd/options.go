package etcd

import (
	"context"
	"time"
)

type Option func(o *options)

type options struct {
	ctx       context.Context
	namespace string
	ttl       time.Duration
	maxRetry  int
}

func WithContext(ctx context.Context) Option {
	return func(o *options) {
		o.ctx = ctx
	}
}

func WithNamespace(namespace string) Option {
	return func(o *options) {
		o.namespace = namespace
	}
}

func WithTTL(ttl time.Duration) Option {
	return func(o *options) {
		o.ttl = ttl
	}
}

func WithMaxRetry(num int) Option {
	return func(o *options) {
		o.maxRetry = num
	}
}
