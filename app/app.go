package app

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/gotechbook/pkg/endpoint"
	"github.com/gotechbook/pkg/logger"
	"github.com/gotechbook/pkg/transport"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Info interface {
	ID() string
	Name() string
	Version() string
	Metadata() map[string]string
	Endpoint() []string
}

type App struct {
	opts     options
	ctx      context.Context
	cancel   func()
	mu       sync.Mutex
	instance *endpoint.Instance
}

func New(opts ...Option) *App {
	o := options{
		ctx:              context.Background(),
		sigs:             []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
		registrarTimeout: 10 * time.Second,
		stopTimeout:      10 * time.Second,
	}
	if id, err := uuid.NewUUID(); err == nil {
		o.id = id.String()
	}
	for _, opt := range opts {
		opt(&o)
	}
	if o.logger != nil {
		logger.SetLogger(o.logger)
	}
	ctx, cancel := context.WithCancel(o.ctx)
	return &App{
		ctx:    ctx,
		cancel: cancel,
		opts:   o,
	}
}

// ID returns app instance id.
func (a *App) ID() string { return a.opts.id }

// Name returns service name.
func (a *App) Name() string { return a.opts.name }

// Version returns app version.
func (a *App) Version() string { return a.opts.version }

// Metadata returns service metadata.
func (a *App) Metadata() map[string]string { return a.opts.metadata }

// Endpoint returns endpoints.
func (a *App) Endpoint() []string {
	if a.instance != nil {
		return a.instance.Endpoints
	}
	return nil
}

func (a *App) beforeStart(ctx context.Context) error {
	for _, fn := range a.opts.beforeStart {
		if err := fn(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (a *App) afterStart(ctx context.Context) error {
	for _, fn := range a.opts.afterStart {
		if err := fn(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (a *App) afterStop(ctx context.Context) error {
	for _, fn := range a.opts.afterStop {
		if err := fn(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (a *App) beforeStop(ctx context.Context) error {
	for _, fn := range a.opts.beforeStop {
		if err := fn(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (a *App) buildInstance() (*endpoint.Instance, error) {
	endpoints := make([]string, 0, len(a.opts.endpoints))
	for _, e := range a.opts.endpoints {
		endpoints = append(endpoints, e.String())
	}

	if len(endpoints) > 0 {
		for _, srv := range a.opts.servers {
			if r, ok := srv.(transport.Endpoint); ok {
				e, err := r.Endpoint()
				if err != nil {
					return nil, err
				}
				endpoints = append(endpoints, e.String())
			}
		}
	}
	return &endpoint.Instance{
		ID:        a.opts.id,
		Name:      a.opts.name,
		Version:   a.opts.version,
		Metadata:  a.opts.metadata,
		Endpoints: endpoints,
	}, nil
}

func (a *App) setInstance(instance *endpoint.Instance) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.instance = instance
}

func (a *App) getInstance() *endpoint.Instance {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.instance == nil {
		return nil
	}
	return a.instance
}

func (a *App) Run() error {
	instance, err := a.buildInstance()
	if err != nil {
		return err
	}
	a.setInstance(instance)
	appContext := NewContext(a.ctx, a)
	if err = a.beforeStart(appContext); err != nil {
		return err
	}
	eg, ctx := errgroup.WithContext(appContext)
	for _, srv := range a.opts.servers {
		srv := srv
		eg.Go(func() error {
			<-ctx.Done() // wait for stop signal
			stopCtx, cancel := context.WithTimeout(NewContext(a.opts.ctx, a), a.opts.stopTimeout)
			defer cancel()
			return srv.Stop(stopCtx)
		})
		eg.Go(func() error {
			return srv.Start(appContext)
		})
	}
	if a.opts.registrar != nil {
		registrarContext, registrarCancel := context.WithTimeout(ctx, a.opts.registrarTimeout)
		defer registrarCancel()
		if err = a.opts.registrar.Register(registrarContext, instance); err != nil {
			return err
		}
	}
	if err = a.afterStart(appContext); err != nil {
		return err
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, a.opts.sigs...)
	eg.Go(func() error {
		select {
		case <-ctx.Done():
			return nil
		case <-c:
			return a.Stop()
		}
	})
	if err = eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	return a.afterStop(appContext)
}

func (a *App) Stop() (err error) {
	ctx := NewContext(a.ctx, a)
	if err = a.beforeStop(ctx); err != nil {
		return err
	}
	instance := a.getInstance()
	if a.opts.registrar != nil && instance != nil {
		ctx, cancel := context.WithTimeout(NewContext(a.ctx, a), a.opts.registrarTimeout)
		defer cancel()
		if err = a.opts.registrar.Deregister(ctx, instance); err != nil {
			return err
		}
	}
	if a.cancel != nil {
		a.cancel()
	}
	return err
}

type appKey struct{}

// NewContext returns a new Context that carries value.
func NewContext(ctx context.Context, s Info) context.Context {
	return context.WithValue(ctx, appKey{}, s)
}

// FromContext returns the Transport value stored in ctx, if any.
func FromContext(ctx context.Context) (s Info, ok bool) {
	s, ok = ctx.Value(appKey{}).(Info)
	return
}
