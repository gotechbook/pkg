package grpc

import (
	"github.com/gotechbook/pkg/transport"
	"google.golang.org/grpc/metadata"
)

type headerCarrier metadata.MD

func (h headerCarrier) Get(k string) string {
	v := metadata.MD(h).Get(k)
	if len(v) > 0 {
		return v[0]
	}
	return ""
}

func (h headerCarrier) Set(k, v string) {
	metadata.MD(h).Set(k, v)
}

func (h headerCarrier) Keys() []string {
	ks := make([]string, 0, len(h))
	for k := range metadata.MD(h) {
		ks = append(ks, k)
	}
	return ks
}

var _ transport.Transporter = (*Transport)(nil)

type Transport struct {
	endpoint    string
	operation   string
	reqHeader   headerCarrier
	replyHeader headerCarrier
}

func (t *Transport) Kind() transport.Kind {
	return transport.KindGRPC
}

func (t *Transport) Endpoint() string {
	return t.endpoint
}

func (t *Transport) Operation() string {
	return t.operation
}

func (t *Transport) RequestHeader() transport.Header {
	return t.reqHeader
}

func (t *Transport) ReplyHeader() transport.Header {
	return t.replyHeader
}
