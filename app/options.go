package app

import (
	"context"
	"github.com/gotechbook/pkg/endpoint"
	"github.com/gotechbook/pkg/logger"
	"github.com/gotechbook/pkg/transport"
	"net/url"
	"os"
	"time"
)

type Option func(o *options)

type options struct {
	id        string
	name      string
	version   string
	metadata  map[string]string
	endpoints []*url.URL

	ctx  context.Context
	sigs []os.Signal

	logger           logger.Logger
	registrar        endpoint.Registrar
	registrarTimeout time.Duration
	stopTimeout      time.Duration
	servers          []transport.Server
	beforeStart      []func(context.Context) error
	beforeStop       []func(context.Context) error
	afterStart       []func(context.Context) error
	afterStop        []func(context.Context) error
}

func WithID(id string) Option {
	return func(o *options) { o.id = id }
}

func WithName(name string) Option {
	return func(o *options) { o.name = name }
}

func WithVersion(version string) Option {
	return func(o *options) { o.version = version }
}

func WithMetadata(md map[string]string) Option {
	return func(o *options) { o.metadata = md }
}

func WithEndpoints(endpoints ...*url.URL) Option {
	return func(o *options) { o.endpoints = endpoints }
}

func WithContext(ctx context.Context) Option {
	return func(o *options) { o.ctx = ctx }
}

func WithLogger(logger logger.Logger) Option {
	return func(o *options) { o.logger = logger }
}

func WithRegistrar(r endpoint.Registrar) Option {
	return func(o *options) { o.registrar = r }
}

func WithRegistrarTimeout(t time.Duration) Option {
	return func(o *options) { o.registrarTimeout = t }
}

func WithStopTimeout(t time.Duration) Option {
	return func(o *options) { o.stopTimeout = t }
}

func WithServer(srv ...transport.Server) Option {
	return func(o *options) { o.servers = srv }
}

func WithBeforeStart(fn func(context.Context) error) Option {
	return func(o *options) { o.beforeStart = append(o.beforeStart, fn) }
}

func WithBeforeStop(fn func(context.Context) error) Option {
	return func(o *options) { o.beforeStop = append(o.beforeStop, fn) }
}

func WithAfterStart(fn func(context.Context) error) Option {
	return func(o *options) { o.afterStart = append(o.afterStart, fn) }
}

func WithAfterStop(fn func(context.Context) error) Option {
	return func(o *options) { o.afterStop = append(o.afterStop, fn) }
}
