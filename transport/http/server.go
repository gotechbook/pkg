package http

import (
	"context"
	"crypto/tls"
	"github.com/gotechbook/pkg/middleware"
	"github.com/valyala/fasthttp"
	"net"
	"net/url"
	"time"
)

type Server struct {
	*fasthttp.Server
	err        error
	address    string
	network    string
	endpoint   *url.URL
	listener   net.Listener
	timeout    time.Duration
	context    context.Context
	tlsConf    *tls.Config
	middleware middleware.Matcher
}

func NewServer() *Server {
	//fasthttp.Serve()
	return nil
}
