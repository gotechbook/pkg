package transport

import (
	"context"
	"net/url"
)

type Kind string

const (
	KindGRPC Kind = "grpc"
	KindHTTP Kind = "http"
)

func (k Kind) String() string { return string(k) }

type Transporter interface {
	Kind() Kind
	Endpoint() string
	Operation() string
	RequestHeader() Header
	ReplyHeader() Header
}

type Server interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

type Header interface {
	Get(k string) string
	Set(k, v string)
	Keys() []string
}

type Endpoint interface {
	Endpoint() (*url.URL, error)
}

type (
	serverTransportKey struct{}
	clientTransportKey struct{}
)

// NewServerContext returns a new Context that carries value.
func NewServerContext(ctx context.Context, tr Transporter) context.Context {
	return context.WithValue(ctx, serverTransportKey{}, tr)
}

// FromServerContext returns the Transport value stored in ctx, if any.
func FromServerContext(ctx context.Context) (tr Transporter, ok bool) {
	tr, ok = ctx.Value(serverTransportKey{}).(Transporter)
	return
}

// NewClientContext returns a new Context that carries value.
func NewClientContext(ctx context.Context, tr Transporter) context.Context {
	return context.WithValue(ctx, clientTransportKey{}, tr)
}

// FromClientContext returns the Transport value stored in ctx, if any.
func FromClientContext(ctx context.Context) (tr Transporter, ok bool) {
	tr, ok = ctx.Value(clientTransportKey{}).(Transporter)
	return
}
