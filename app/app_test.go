package app

import (
	"context"
	"github.com/gotechbook/pkg/endpoint/etcd"
	"github.com/gotechbook/pkg/logger"
	"github.com/gotechbook/pkg/transport/grpc"
	clientv3 "go.etcd.io/etcd/client/v3"
	"testing"
	"time"
)

func TestApp(t *testing.T) {
	logger.SetLogger(logger.With(logger.DefaultLogger, "caller", logger.DefaultCaller, "ts", logger.DefaultTimestamp))
	gs := grpc.NewServer(grpc.WithServerAddress(":9001"))
	client, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"127.0.0.1:65160"},
	})
	if err != nil {
		logger.Fatal(err)
	}
	r := etcd.New(client)
	app := New(
		WithName("test"),
		WithVersion("v1.0.0"),
		WithServer(gs),
		WithRegistrar(r),
		WithBeforeStart(func(ctx context.Context) error {
			t.Log("BeforeStart...")
			return nil
		}),
		WithBeforeStop(func(ctx context.Context) error {
			t.Log("BeforeStop...")
			return nil
		}),
		WithAfterStart(func(ctx context.Context) error {
			t.Log("AfterStart...")
			return nil
		}),
		WithAfterStop(func(ctx context.Context) error {
			t.Log("AfterStop...")
			return nil
		}),
	)
	time.AfterFunc(time.Second, func() {
		_ = app.Stop()
	})
	if err := app.Run(); err != nil {
		t.Fatal(err)
	}
}
